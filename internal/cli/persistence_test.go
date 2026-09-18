package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestScanStatusPersistsAcrossCLIInvocations(t *testing.T) {
	t.Setenv("KRYPTASEC_DATA_DIR", t.TempDir())

	var createOut, createErr bytes.Buffer
	code := Run(context.Background(), []string{"scan", t.TempDir()}, &createOut, &createErr, "dev")
	if code != 0 {
		t.Fatalf("scan exit code = %d; stderr=%q", code, createErr.String())
	}

	var scanID string
	for _, line := range strings.Split(createOut.String(), "\n") {
		if strings.HasPrefix(line, "scan: ") {
			scanID = strings.TrimSpace(strings.TrimPrefix(line, "scan: "))
		}
	}
	if scanID == "" {
		t.Fatalf("scan ID not found in output: %q", createOut.String())
	}

	var statusOut, statusErr bytes.Buffer
	code = Run(context.Background(), []string{"scan", "status", scanID}, &statusOut, &statusErr, "dev")
	if code != 0 {
		t.Fatalf("status exit code = %d; stderr=%q", code, statusErr.String())
	}
	if !strings.Contains(statusOut.String(), "status: ready") {
		t.Fatalf("unexpected status output: %q", statusOut.String())
	}
	if !strings.Contains(statusOut.String(), "scan: "+scanID) {
		t.Fatalf("status output missing scan ID: %q", statusOut.String())
	}
}

func TestScanStatusUnknownIDReturnsNotFound(t *testing.T) {
	t.Setenv("KRYPTASEC_DATA_DIR", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), []string{"scan", "status", "scan_missing"}, &stdout, &stderr, "dev")
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "scan not found") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}
