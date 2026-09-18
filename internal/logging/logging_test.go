package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewRedactsSensitiveStructuredFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, slog.LevelDebug)

	logger.Info(
		"request",
		"scan_id", "scan_123",
		"token", "super-secret-token",
		"authorization", "Bearer abc",
		"cookie", "session=secret",
		"password", "hunter2",
	)

	out := buf.String()
	if strings.Contains(out, "super-secret-token") ||
		strings.Contains(out, "Bearer abc") ||
		strings.Contains(out, "session=secret") ||
		strings.Contains(out, "hunter2") {
		t.Fatalf("sensitive value leaked in log: %s", out)
	}
	if !strings.Contains(out, "scan_123") {
		t.Fatalf("non-sensitive field missing from log: %s", out)
	}
	if strings.Count(out, "[REDACTED]") < 4 {
		t.Fatalf("expected sensitive fields to be redacted: %s", out)
	}
}

func TestNewHonorsConfiguredLogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, slog.LevelWarn)

	logger.Info("info-hidden")
	if buf.Len() != 0 {
		t.Fatalf("info log should be suppressed at warn level: %s", buf.String())
	}

	logger.Warn("warn-visible")
	if !strings.Contains(buf.String(), "warn-visible") {
		t.Fatalf("warn log missing: %s", buf.String())
	}
}
