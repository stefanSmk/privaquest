package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/stefanSmk/privaquest/internal/i18n"
	"github.com/stefanSmk/privaquest/internal/store"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidType  = errors.New("invalid request type")
	ErrInvalidLocale = errors.New("unsupported locale")
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type Service struct {
	store *store.Store
}

func New(st *store.Store) *Service {
	return &Service{store: st}
}

type SubmitInput struct {
	Type        store.RequestType
	Email       string
	FullName    string
	Description string
	Locale      string
}

type SubmitResult struct {
	Reference string    `json:"reference"`
	Message   string    `json:"message"`
	DueAt     time.Time `json:"due_at"`
	Locale    string    `json:"locale"`
}

type StatusUpdate struct {
	Status store.RequestStatus
	Notes  string
	Actor  string
}

func (s *Service) Submit(ctx context.Context, in SubmitInput) (SubmitResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if !emailPattern.MatchString(email) {
		return SubmitResult{}, ErrInvalidEmail
	}
	if !validType(in.Type) {
		return SubmitResult{}, ErrInvalidType
	}
	locale := i18n.Normalize(in.Locale)
	if locale != "en" && locale != "de" && locale != "fr" {
		return SubmitResult{}, ErrInvalidLocale
	}

	ref, err := newReference()
	if err != nil {
		return SubmitResult{}, err
	}

	due := time.Now().UTC().Add(30 * 24 * time.Hour)
	req, err := s.store.CreateRequest(ctx, store.PrivacyRequest{
		Reference:   ref,
		Type:        in.Type,
		Email:       email,
		FullName:    strings.TrimSpace(in.FullName),
		Description: strings.TrimSpace(in.Description),
		Locale:      locale,
		Status:      store.StatusOpen,
		DueAt:       due,
	})
	if err != nil {
		return SubmitResult{}, err
	}

	_ = s.store.AddAudit(ctx, store.AuditEntry{
		RequestID: req.ID,
		Action:    "request.created",
		Actor:     "public",
		Details:   fmt.Sprintf("type=%s locale=%s", req.Type, req.Locale),
	})

	return SubmitResult{
		Reference: ref,
		Message:   i18n.T(locale, "request_received"),
		DueAt:     due,
		Locale:    locale,
	}, nil
}

func (s *Service) GetPublicStatus(ctx context.Context, reference, email string) (store.PrivacyRequest, []store.AuditEntry, error) {
	req, err := s.store.GetByReference(ctx, reference)
	if err != nil {
		return store.PrivacyRequest{}, nil, err
	}
	if strings.ToLower(strings.TrimSpace(email)) != req.Email {
		return store.PrivacyRequest{}, nil, store.ErrNotFound
	}
	trail, err := s.store.AuditTrail(ctx, req.ID)
	return req, trail, err
}

func (s *Service) List(ctx context.Context, status store.RequestStatus) ([]store.PrivacyRequest, error) {
	return s.store.List(ctx, status, 100)
}

func (s *Service) Dashboard(ctx context.Context) (store.Dashboard, error) {
	return s.store.Dashboard(ctx)
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, upd StatusUpdate) (store.PrivacyRequest, error) {
	if upd.Actor == "" {
		upd.Actor = "admin"
	}
	if err := s.store.UpdateStatus(ctx, id, upd.Status, upd.Notes); err != nil {
		return store.PrivacyRequest{}, err
	}
	_ = s.store.AddAudit(ctx, store.AuditEntry{
		RequestID: id,
		Action:    "request.status_changed",
		Actor:     upd.Actor,
		Details:   fmt.Sprintf("status=%s notes=%q", upd.Status, upd.Notes),
	})
	return s.store.GetByID(ctx, id)
}

func (s *Service) AuditTrail(ctx context.Context, id int64) ([]store.AuditEntry, error) {
	return s.store.AuditTrail(ctx, id)
}

func validType(t store.RequestType) bool {
	switch t {
	case store.TypeAccess, store.TypeDelete, store.TypeRectify, store.TypeObject:
		return true
	default:
		return false
	}
}

func newReference() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	year := time.Now().UTC().Year()
	return fmt.Sprintf("DSR-%d-%s", year, strings.ToUpper(hex.EncodeToString(b))), nil
}
