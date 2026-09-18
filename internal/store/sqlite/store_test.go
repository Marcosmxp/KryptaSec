package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/scan"
	"github.com/Marcosmxp/KryptaSec/internal/store"
)

func TestStoreCreateGetAndReopen(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "kryptasec.db")

	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	now := time.Date(2026, 9, 18, 20, 0, 0, 0, time.UTC)
	job, err := scan.New("local", "/tmp/app", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Create(ctx, job); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen error = %v", err)
	}
	t.Cleanup(func() { _ = reopened.Close() })

	got, err := reopened.Get(ctx, job.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != job.ID || got.Status != scan.StatusCreated || got.Target != job.Target {
		t.Fatalf("unexpected scan after reopen: %+v", got)
	}
}

func TestStoreTransitionIsAtomicAndChecksExpectedState(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "kryptasec.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	now := time.Now().UTC()
	job, _ := scan.New("http", "https://example.com", now)
	if err := s.Create(ctx, job); err != nil {
		t.Fatal(err)
	}

	updated, err := s.Transition(ctx, job.ID, scan.StatusCreated, scan.StatusValidatingScope, now.Add(time.Second))
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	if updated.Status != scan.StatusValidatingScope {
		t.Fatalf("Status = %q, want validating_scope", updated.Status)
	}

	if _, err := s.Transition(ctx, job.ID, scan.StatusCreated, scan.StatusFailed, now.Add(2*time.Second)); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	got, err := s.Get(ctx, job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != scan.StatusValidatingScope {
		t.Fatalf("failed transition changed persisted status to %q", got.Status)
	}
}

func TestStoreGetUnknownScanReturnsNotFound(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "kryptasec.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	if _, err := s.Get(context.Background(), "scan_missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestOpenSetsCurrentSchemaVersion(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "kryptasec.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	version, err := s.SchemaVersion(context.Background())
	if err != nil {
		t.Fatalf("SchemaVersion() error = %v", err)
	}
	if version != CurrentSchemaVersion {
		t.Fatalf("schema version = %d, want %d", version, CurrentSchemaVersion)
	}
	if version != 1 {
		t.Fatalf("initial schema version = %d, want 1", version)
	}
}

func TestOpenRejectsDatabaseFromNewerSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "future.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("PRAGMA user_version = 2"); err != nil {
		_ = db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Open(path)
	if err == nil {
		t.Fatal("Open() expected error for future schema")
	}
	if !strings.Contains(err.Error(), "newer schema version") {
		t.Fatalf("unexpected error: %v", err)
	}
}
