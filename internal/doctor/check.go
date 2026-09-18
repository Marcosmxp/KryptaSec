package doctor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Marcosmxp/KryptaSec/internal/config"
	sqlitestore "github.com/Marcosmxp/KryptaSec/internal/store/sqlite"
)

type Check struct {
	Name    string
	OK      bool
	Message string
}

type Report struct {
	Checks []Check
}

func (r Report) OK() bool {
	for _, check := range r.Checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func Run(ctx context.Context, cfg config.Config) Report {
	checks := []Check{{
		Name:    "runtime",
		OK:      ctx.Err() == nil,
		Message: fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}}

	dataDir := checkDataDir(cfg.DataDir)
	checks = append(checks, dataDir)
	if dataDir.OK {
		checks = append(checks, checkDatabase(ctx, cfg.DataDir))
	}

	return Report{Checks: checks}
}

func checkDataDir(dir string) Check {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return Check{Name: "data_dir", OK: false, Message: err.Error()}
	}

	f, err := os.CreateTemp(dir, ".doctor-*")
	if err != nil {
		return Check{Name: "data_dir", OK: false, Message: err.Error()}
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		return Check{Name: "data_dir", OK: false, Message: err.Error()}
	}
	if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
		return Check{Name: "data_dir", OK: false, Message: err.Error()}
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	return Check{Name: "data_dir", OK: true, Message: abs}
}

func checkDatabase(ctx context.Context, dataDir string) Check {
	path := filepath.Join(dataDir, "kryptasec.db")
	s, err := sqlitestore.Open(path)
	if err != nil {
		return Check{Name: "database", OK: false, Message: err.Error()}
	}
	defer func() { _ = s.Close() }()

	if err := s.Health(ctx); err != nil {
		return Check{Name: "database", OK: false, Message: err.Error()}
	}

	version, err := s.SchemaVersion(ctx)
	if err != nil {
		return Check{Name: "database", OK: false, Message: err.Error()}
	}
	return Check{
		Name:    "database",
		OK:      true,
		Message: fmt.Sprintf("healthy schema=%d", version),
	}
}
