package app

import (
	"bytes"
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	projectlogging "github.com/Marcosmxp/KryptaSec/internal/logging"
	sqlitestore "github.com/Marcosmxp/KryptaSec/internal/store/sqlite"
)

func TestStartScanLogsLifecycleWithoutTargetValue(t *testing.T) {
	var logs bytes.Buffer
	logger := projectlogging.New(&logs, slog.LevelDebug)

	s, err := sqlitestore.Open(filepath.Join(t.TempDir(), "kryptasec.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })

	targetDir := t.TempDir()
	svc := Service{
		Store:  s,
		Logger: logger,
		Now:    func() time.Time { return time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC) },
	}

	result, err := svc.StartScan(context.Background(), StartScanRequest{Target: targetDir})
	if err != nil {
		t.Fatalf("StartScan() error = %v", err)
	}

	out := logs.String()
	if !strings.Contains(out, result.Scan.ID) {
		t.Fatalf("structured log missing scan ID: %s", out)
	}
	if !strings.Contains(out, "ready") {
		t.Fatalf("structured log missing final status: %s", out)
	}
	if strings.Contains(out, result.Scan.Target) {
		t.Fatalf("target value must not be written to lifecycle logs: %s", out)
	}
}
