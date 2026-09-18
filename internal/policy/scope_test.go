package policy

import (
	"testing"

	"github.com/Marcosmxp/KryptaSec/internal/target"
)

func TestScopeAllowsLocalTarget(t *testing.T) {
	if !(Scope{}).Allows(target.Target{Kind: target.KindLocal, Canonical: "/tmp/app"}) {
		t.Fatal("local target should be allowed")
	}
}

func TestScopeAllowsExactRemoteHost(t *testing.T) {
	s := Scope{Hosts: []string{"example.com"}}
	if !s.Allows(target.Target{Kind: target.KindHTTP, Canonical: "https://example.com/app"}) {
		t.Fatal("exact host should be allowed")
	}
}

func TestScopeHostComparisonIsCaseInsensitive(t *testing.T) {
	s := Scope{Hosts: []string{"EXAMPLE.COM"}}
	if !s.Allows(target.Target{Kind: target.KindHTTP, Canonical: "https://example.com"}) {
		t.Fatal("host comparison should be case-insensitive")
	}
}

func TestScopeDeniesUnlistedSubdomain(t *testing.T) {
	s := Scope{Hosts: []string{"example.com"}}
	if s.Allows(target.Target{Kind: target.KindHTTP, Canonical: "https://api.example.com"}) {
		t.Fatal("subdomain must not be implicitly authorized")
	}
}

func TestScopeDoesNotExpandWildcard(t *testing.T) {
	s := Scope{Hosts: []string{"*.example.com"}}
	if s.Allows(target.Target{Kind: target.KindHTTP, Canonical: "https://api.example.com"}) {
		t.Fatal("wildcards must not be expanded in Phase 1")
	}
}
