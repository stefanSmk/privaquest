package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS privacy_requests (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	reference TEXT NOT NULL UNIQUE,
	type TEXT NOT NULL,
	email TEXT NOT NULL,
	full_name TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	locale TEXT NOT NULL DEFAULT 'en',
	status TEXT NOT NULL DEFAULT 'open',
	due_at TEXT NOT NULL,
	created_at TEXT NOT NULL,
	resolved_at TEXT,
	notes TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS audit_log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	request_id INTEGER NOT NULL,
	action TEXT NOT NULL,
	actor TEXT NOT NULL DEFAULT 'system',
	details TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL,
	FOREIGN KEY(request_id) REFERENCES privacy_requests(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_requests_status ON privacy_requests(status);
CREATE INDEX IF NOT EXISTS idx_requests_due_at ON privacy_requests(due_at);
`)
	return err
}

func (s *Store) CreateRequest(ctx context.Context, req PrivacyRequest) (PrivacyRequest, error) {
	now := time.Now().UTC()
	req.CreatedAt = now
	if req.DueAt.IsZero() {
		req.DueAt = now.Add(30 * 24 * time.Hour)
	}
	if req.Status == "" {
		req.Status = StatusOpen
	}

	res, err := s.db.ExecContext(ctx, `
INSERT INTO privacy_requests (reference, type, email, full_name, description, locale, status, due_at, created_at, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.Reference, req.Type, req.Email, req.FullName, req.Description, req.Locale,
		req.Status, req.DueAt.Format(time.RFC3339), req.CreatedAt.Format(time.RFC3339), req.Notes,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return PrivacyRequest{}, errors.New("duplicate reference")
		}
		return PrivacyRequest{}, err
	}
	id, _ := res.LastInsertId()
	req.ID = id
	return req, nil
}

func (s *Store) GetByReference(ctx context.Context, ref string) (PrivacyRequest, error) {
	return s.scanRequest(s.db.QueryRowContext(ctx, `
SELECT id, reference, type, email, full_name, description, locale, status, due_at, created_at, resolved_at, notes
FROM privacy_requests WHERE reference = ?`, ref))
}

func (s *Store) GetByID(ctx context.Context, id int64) (PrivacyRequest, error) {
	return s.scanRequest(s.db.QueryRowContext(ctx, `
SELECT id, reference, type, email, full_name, description, locale, status, due_at, created_at, resolved_at, notes
FROM privacy_requests WHERE id = ?`, id))
}

func (s *Store) List(ctx context.Context, status RequestStatus, limit int) ([]PrivacyRequest, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `
SELECT id, reference, type, email, full_name, description, locale, status, due_at, created_at, resolved_at, notes
FROM privacy_requests`
	args := []any{}
	if status != "" {
		query += ` WHERE status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY due_at ASC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PrivacyRequest
	for rows.Next() {
		req, err := s.scanRequestRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, req)
	}
	return out, rows.Err()
}

func (s *Store) UpdateStatus(ctx context.Context, id int64, status RequestStatus, notes string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	var resolved any
	if status == StatusResolved || status == StatusRejected {
		resolved = now
	}
	res, err := s.db.ExecContext(ctx, `
UPDATE privacy_requests SET status = ?, notes = ?, resolved_at = COALESCE(?, resolved_at) WHERE id = ?`,
		status, notes, resolved, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AddAudit(ctx context.Context, entry AuditEntry) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO audit_log (request_id, action, actor, details, created_at) VALUES (?, ?, ?, ?, ?)`,
		entry.RequestID, entry.Action, entry.Actor, entry.Details, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

func (s *Store) AuditTrail(ctx context.Context, requestID int64) ([]AuditEntry, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, request_id, action, actor, details, created_at FROM audit_log
WHERE request_id = ? ORDER BY created_at ASC`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditEntry
	for rows.Next() {
		var e AuditEntry
		var created string
		if err := rows.Scan(&e.ID, &e.RequestID, &e.Action, &e.Actor, &e.Details, &created); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) Dashboard(ctx context.Context) (Dashboard, error) {
	var d Dashboard
	now := time.Now().UTC()
	in7 := now.Add(7 * 24 * time.Hour).Format(time.RFC3339)
	last30 := now.Add(-30 * 24 * time.Hour).Format(time.RFC3339)

	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM privacy_requests WHERE status IN ('open','in_progress')`).Scan(&d.OpenCount)
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM privacy_requests WHERE status IN ('open','in_progress') AND due_at < ?`, now.Format(time.RFC3339)).Scan(&d.OverdueCount)
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM privacy_requests WHERE status IN ('open','in_progress') AND due_at <= ?`, in7).Scan(&d.DueWithin7Days)
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM privacy_requests WHERE status = 'resolved' AND resolved_at >= ?`, last30).Scan(&d.ResolvedLast30d)
	return d, nil
}

func (s *Store) scanRequest(row *sql.Row) (PrivacyRequest, error) {
	var req PrivacyRequest
	var due, created string
	var resolved sql.NullString
	err := row.Scan(&req.ID, &req.Reference, &req.Type, &req.Email, &req.FullName, &req.Description,
		&req.Locale, &req.Status, &due, &created, &resolved, &req.Notes)
	if errors.Is(err, sql.ErrNoRows) {
		return PrivacyRequest{}, ErrNotFound
	}
	if err != nil {
		return PrivacyRequest{}, err
	}
	req.DueAt, _ = time.Parse(time.RFC3339, due)
	req.CreatedAt, _ = time.Parse(time.RFC3339, created)
	if resolved.Valid {
		t, _ := time.Parse(time.RFC3339, resolved.String)
		req.ResolvedAt = &t
	}
	return req, nil
}

func (s *Store) scanRequestRows(rows *sql.Rows) (PrivacyRequest, error) {
	var req PrivacyRequest
	var due, created string
	var resolved sql.NullString
	err := rows.Scan(&req.ID, &req.Reference, &req.Type, &req.Email, &req.FullName, &req.Description,
		&req.Locale, &req.Status, &due, &created, &resolved, &req.Notes)
	if err != nil {
		return PrivacyRequest{}, err
	}
	req.DueAt, _ = time.Parse(time.RFC3339, due)
	req.CreatedAt, _ = time.Parse(time.RFC3339, created)
	if resolved.Valid {
		t, _ := time.Parse(time.RFC3339, resolved.String)
		req.ResolvedAt = &t
	}
	return req, nil
}
