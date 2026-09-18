package policy

import (
	"net/url"
	"strings"

	"github.com/Marcosmxp/KryptaSec/internal/target"
)

type Scope struct {
	Hosts []string
}

func (s Scope) Allows(t target.Target) bool {
	if t.Kind == target.KindLocal {
		return true
	}
	if t.Kind != target.KindHTTP {
		return false
	}

	u, err := url.Parse(t.Canonical)
	if err != nil {
		return false
	}
	targetHost := strings.ToLower(u.Hostname())
	if targetHost == "" {
		return false
	}

	for _, allowed := range s.Hosts {
		if strings.EqualFold(strings.TrimSpace(allowed), targetHost) {
			return true
		}
	}
	return false
}
