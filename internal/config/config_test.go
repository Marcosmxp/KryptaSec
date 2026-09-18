package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("KRYPTASEC_DATA_DIR", dir)
	t.Setenv("KRYPTASEC_LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DataDir != dir {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, dir)
	}
	if cfg.LogLevel != slog.LevelDebug {
		t.Fatalf("LogLevel = %v, want debug", cfg.LogLevel)
	}
}

func TestLoadUsesUserConfigDirectoryByDefault(t *testing.T) {
	t.Setenv("KRYPTASEC_DATA_DIR", "")
	t.Setenv("KRYPTASEC_LOG_LEVEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	base, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir() error = %v", err)
	}
	want := filepath.Join(base, "kryptasec")
	if cfg.DataDir != want {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, want)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Fatalf("LogLevel = %v, want info", cfg.LogLevel)
	}
}

func TestLoadRejectsInvalidLogLevel(t *testing.T) {
	t.Setenv("KRYPTASEC_LOG_LEVEL", "verbose")
	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error for invalid log level")
	}
}
