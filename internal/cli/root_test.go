package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(context.Background(), []string{"--version"}, &stdout, &stderr, "0.1.0-dev")

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if got := strings.TrimSpace(stdout.String()); got != "KryptaSec 0.1.0-dev" {
		t.Fatalf("unexpected version output: %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunUnknownCommandFailsClosed(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(context.Background(), []string{"attack"}, &stdout, &stderr, "dev")

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected unknown command error, got %q", stderr.String())
	}
}

func TestRunDoctorReportsHealthyChecks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	dataDir := t.TempDir()
	t.Setenv("KRYPTASEC_DATA_DIR", dataDir)

	code := Run(context.Background(), []string{"doctor"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "[ok] runtime") || !strings.Contains(output, "[ok] data_dir") {
		t.Fatalf("unexpected doctor output: %q", output)
	}
}

func TestRunScanLocalTarget(t *testing.T) {
	var stdout, stderr bytes.Buffer
	dir := t.TempDir()

	code := Run(context.Background(), []string{"scan", dir}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "status: ready") || !strings.Contains(stdout.String(), "kind: local") {
		t.Fatalf("unexpected scan output: %q", stdout.String())
	}
}

func TestRunScanRemoteRequiresScopeHost(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(context.Background(), []string{"scan", "https://example.com"}, &stdout, &stderr, "dev")

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "outside the authorized scope") {
		t.Fatalf("unexpected error: %q", stderr.String())
	}
}

func TestRunScanRemoteAcceptsExactScopeHost(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(context.Background(), []string{"scan", "--scope-host", "example.com", "https://example.com"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "status: ready") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}
