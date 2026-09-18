package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeLocalDirectory(t *testing.T) {
	dir := t.TempDir()
	got, err := Normalize(dir)
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	abs, _ := filepath.Abs(dir)
	if got.Kind != KindLocal || got.Canonical != filepath.Clean(abs) {
		t.Fatalf("unexpected target: %+v", got)
	}
}

func TestNormalizeRejectsMissingLocalPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	if _, err := Normalize(path); err == nil {
		t.Fatal("Normalize() expected error for missing path")
	}
}

func TestNormalizeHTTPURL(t *testing.T) {
	got, err := Normalize("https://EXAMPLE.com:443/a")
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if got.Kind != KindHTTP || got.Canonical != "https://example.com/a" {
		t.Fatalf("unexpected target: %+v", got)
	}
}

func TestNormalizeRejectsUnsupportedScheme(t *testing.T) {
	if _, err := Normalize("ftp://example.com/file"); err == nil {
		t.Fatal("Normalize() expected error for unsupported scheme")
	}
}

func TestNormalizeRejectsURLCredentials(t *testing.T) {
	if _, err := Normalize("https://user:pass@example.com"); err == nil {
		t.Fatal("Normalize() expected error for URL credentials")
	}
}

func TestNormalizeRejectsRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Normalize(path); err == nil {
		t.Fatal("Normalize() expected error for regular file")
	}
}
