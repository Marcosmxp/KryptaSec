package target

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Kind string

const (
	KindLocal Kind = "local"
	KindHTTP  Kind = "http"
)

type Target struct {
	Kind      Kind
	Raw       string
	Canonical string
}

func Normalize(raw string) (Target, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Target{}, fmt.Errorf("target is required")
	}

	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return normalizeHTTP(raw)
	}
	if strings.Contains(raw, "://") {
		u, _ := url.Parse(raw)
		scheme := "unknown"
		if u != nil && u.Scheme != "" {
			scheme = u.Scheme
		}
		return Target{}, fmt.Errorf("unsupported target scheme %q", scheme)
	}

	return normalizeLocal(raw)
}

func normalizeLocal(raw string) (Target, error) {
	abs, err := filepath.Abs(raw)
	if err != nil {
		return Target{}, fmt.Errorf("resolve target path: %w", err)
	}
	abs = filepath.Clean(abs)
	info, err := os.Stat(abs)
	if err != nil {
		return Target{}, fmt.Errorf("stat target path: %w", err)
	}
	if !info.IsDir() {
		return Target{}, fmt.Errorf("local target must be a directory")
	}
	return Target{Kind: KindLocal, Raw: raw, Canonical: abs}, nil
}

func normalizeHTTP(raw string) (Target, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return Target{}, fmt.Errorf("parse target URL: %w", err)
	}
	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme != "http" && u.Scheme != "https" {
		return Target{}, fmt.Errorf("unsupported target scheme %q", u.Scheme)
	}
	if u.User != nil {
		return Target{}, fmt.Errorf("URL credentials are not allowed")
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return Target{}, fmt.Errorf("target URL host is required")
	}

	port := u.Port()
	if (u.Scheme == "https" && port == "443") || (u.Scheme == "http" && port == "80") {
		port = ""
	}
	if port == "" {
		if strings.Contains(host, ":") {
			u.Host = "[" + host + "]"
		} else {
			u.Host = host
		}
	} else {
		u.Host = net.JoinHostPort(host, port)
	}

	return Target{Kind: KindHTTP, Raw: raw, Canonical: u.String()}, nil
}
