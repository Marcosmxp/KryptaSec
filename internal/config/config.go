package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DataDir  string
	LogLevel slog.Level
}

func Load() (Config, error) {
	dataDir := strings.TrimSpace(os.Getenv("KRYPTASEC_DATA_DIR"))
	if dataDir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return Config{}, fmt.Errorf("resolve user config directory: %w", err)
		}
		dataDir = filepath.Join(base, "kryptasec")
	}

	level, err := parseLogLevel(os.Getenv("KRYPTASEC_LOG_LEVEL"))
	if err != nil {
		return Config{}, err
	}

	return Config{DataDir: filepath.Clean(dataDir), LogLevel: level}, nil
}

func parseLogLevel(raw string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid KRYPTASEC_LOG_LEVEL %q: expected debug, info, warn, or error", raw)
	}
}
