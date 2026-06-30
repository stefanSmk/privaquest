package store

import "time"

type RequestType string

const (
	TypeAccess  RequestType = "access"
	TypeDelete  RequestType = "delete"
	TypeRectify RequestType = "rectify"
	TypeObject  RequestType = "object"
)

type RequestStatus string

const (
	StatusOpen       RequestStatus = "open"
	StatusInProgress RequestStatus = "in_progress"
	StatusResolved   RequestStatus = "resolved"
	StatusRejected   RequestStatus = "rejected"
)

type PrivacyRequest struct {
	ID          int64
	Reference   string
	Type        RequestType
	Email       string
	FullName    string
	Description string
	Locale      string
	Status      RequestStatus
	DueAt       time.Time
	CreatedAt   time.Time
	ResolvedAt  *time.Time
	Notes       string
}

type AuditEntry struct {
	ID        int64
	RequestID int64
	Action    string
	Actor     string
	Details   string
	CreatedAt time.Time
}

type Dashboard struct {
	OpenCount       int64
	OverdueCount    int64
	DueWithin7Days  int64
	ResolvedLast30d int64
}
