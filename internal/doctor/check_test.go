package doctor

import (
	"context"
	"strings"
	"testing"

	"github.com/Marcosmxp/KryptaSec/internal/config"
)

func TestRunReportsWritableDataDirectoryAndDatabaseHealth(t *testing.T) {
	report := Run(context.Background(), config.Config{DataDir: t.TempDir()})

	if !report.OK() {
		t.Fatalf("expected healthy report: %+v", report.Checks)
	}

	checks := map[string]Check{}
	for _, check := range report.Checks {
		checks[check.Name] = check
	}

	if !checks["runtime"].OK {
		t.Fatalf("runtime check missing or unhealthy: %+v", checks)
	}
	if !checks["data_dir"].OK {
		t.Fatalf("data_dir check missing or unhealthy: %+v", checks)
	}
	if !checks["database"].OK {
		t.Fatalf("database check missing or unhealthy: %+v", checks)
	}
	if !strings.Contains(checks["database"].Message, "schema=1") {
		t.Fatalf("database check should report schema version: %+v", checks["database"])
	}
}
