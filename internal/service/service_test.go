package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stefanSmk/privaquest/internal/store"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open("file:" + filepath.Join(t.TempDir(), "test.db") + "?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	st := store.New(db)
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return New(st)
}

func TestSubmitRequest(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	res, err := svc.Submit(ctx, SubmitInput{
		Type:        store.TypeAccess,
		Email:       "user@example.com",
		FullName:    "Anna",
		Description: "Please send my data",
		Locale:      "de",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Reference == "" || res.Locale != "de" {
		t.Fatalf("unexpected result: %+v", res)
	}
	if res.DueAt.Sub(time.Now().UTC()) < 29*24*time.Hour {
		t.Fatalf("due date should be ~30 days out, got %v", res.DueAt)
	}
}

func TestSubmitValidation(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Submit(context.Background(), SubmitInput{
		Type:  store.TypeDelete,
		Email: "not-an-email",
	})
	if err != ErrInvalidEmail {
		t.Fatalf("expected invalid email, got %v", err)
	}
}

func TestPublicStatus(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	res, err := svc.Submit(ctx, SubmitInput{
		Type:   store.TypeDelete,
		Email:  "test@example.com",
		Locale: "fr",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, _, err := svc.GetPublicStatus(ctx, res.Reference, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if req.Locale != "fr" {
		t.Fatalf("expected fr locale, got %s", req.Locale)
	}

	_, _, err = svc.GetPublicStatus(ctx, res.Reference, "wrong@example.com")
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for wrong email")
	}
}

func TestAdminWorkflow(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	res, _ := svc.Submit(ctx, SubmitInput{
		Type:  store.TypeAccess,
		Email: "admin-test@example.com",
		Locale: "en",
	})
	req, _, _ := svc.GetPublicStatus(ctx, res.Reference, "admin-test@example.com")

	updated, err := svc.UpdateStatus(ctx, req.ID, StatusUpdate{
		Status: store.StatusInProgress,
		Notes:  "Verifying identity",
		Actor:  "admin",
	})
	if err != nil || updated.Status != store.StatusInProgress {
		t.Fatalf("update failed: %+v err=%v", updated, err)
	}

	trail, err := svc.AuditTrail(ctx, req.ID)
	if err != nil || len(trail) < 2 {
		t.Fatalf("expected audit trail, got %d entries", len(trail))
	}
}
