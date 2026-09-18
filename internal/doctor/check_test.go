package doctor

import (
	"context"
	"testing"

	"github.com/Marcosmxp/KryptaSec/internal/config"
)

func TestRunReportsWritableDataDirectory(t *testing.T) {
	report := Run(context.Background(), config.Config{DataDir: t.TempDir()})

	if !report.OK() {
		t.Fatalf("expected healthy report: %+v", report.Checks)
	}
	if len(report.Checks) < 2 {
		t.Fatalf("expected at least runtime and data directory checks, got %d", len(report.Checks))
	}
}
