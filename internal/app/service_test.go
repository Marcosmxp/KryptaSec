package app

import (
	"context"
	"testing"
	"time"
)

func TestStartScanAllowsLocalDirectory(t *testing.T) {
	svc := Service{Now: func() time.Time { return time.Date(2026, 9, 18, 20, 0, 0, 0, time.UTC) }}

	got, err := svc.StartScan(context.Background(), StartScanRequest{Target: t.TempDir()})
	if err != nil {
		t.Fatalf("StartScan() error = %v", err)
	}
	if got.Scan.Status != "ready" || got.Scan.TargetKind != "local" {
		t.Fatalf("unexpected scan: %+v", got.Scan)
	}
}

func TestStartScanDeniesRemoteTargetWithoutExplicitScope(t *testing.T) {
	svc := Service{Now: time.Now}

	got, err := svc.StartScan(context.Background(), StartScanRequest{Target: "https://example.com"})
	if err == nil {
		t.Fatal("StartScan() expected scope error")
	}
	if got.Scan.Status != "failed" {
		t.Fatalf("Status = %q, want failed", got.Scan.Status)
	}
}

func TestStartScanAllowsExactRemoteScopeWithoutNetworkRequest(t *testing.T) {
	svc := Service{Now: time.Now}

	got, err := svc.StartScan(context.Background(), StartScanRequest{
		Target:     "https://EXAMPLE.com:443/app",
		ScopeHosts: []string{"example.com"},
	})
	if err != nil {
		t.Fatalf("StartScan() error = %v", err)
	}
	if got.Scan.Status != "ready" || got.Scan.Target != "https://example.com/app" {
		t.Fatalf("unexpected scan: %+v", got.Scan)
	}
}
