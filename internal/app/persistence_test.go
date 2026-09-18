package app

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/scan"
	sqlitestore "github.com/Marcosmxp/KryptaSec/internal/store/sqlite"
)

func TestStartScanPersistsEveryLifecycleTransition(t *testing.T) {
	ctx := context.Background()
	s, err := sqlitestore.Open(filepath.Join(t.TempDir(), "kryptasec.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	now := time.Date(2026, 9, 18, 20, 0, 0, 0, time.UTC)
	svc := Service{
		Store: s,
		Now:   func() time.Time { return now },
	}

	result, err := svc.StartScan(ctx, StartScanRequest{Target: t.TempDir()})
	if err != nil {
		t.Fatalf("StartScan() error = %v", err)
	}

	persisted, err := s.Get(ctx, result.Scan.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if persisted.Status != scan.StatusReady {
		t.Fatalf("persisted Status = %q, want ready", persisted.Status)
	}
	if !persisted.UpdatedAt.Equal(now) {
		t.Fatalf("UpdatedAt = %v, want %v", persisted.UpdatedAt, now)
	}
}

func TestGetScanReadsPersistedJob(t *testing.T) {
	ctx := context.Background()
	s, err := sqlitestore.Open(filepath.Join(t.TempDir(), "kryptasec.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	now := time.Now().UTC()
	job, _ := scan.New("local", "/tmp/app", now)
	if err := s.Create(ctx, job); err != nil {
		t.Fatal(err)
	}

	svc := Service{Store: s}
	got, err := svc.GetScan(ctx, job.ID)
	if err != nil {
		t.Fatalf("GetScan() error = %v", err)
	}
	if got.ID != job.ID {
		t.Fatalf("ID = %q, want %q", got.ID, job.ID)
	}
}
