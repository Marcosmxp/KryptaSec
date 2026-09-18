package doctor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Marcosmxp/KryptaSec/internal/config"
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

	checks = append(checks, checkDataDir(cfg.DataDir))
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
