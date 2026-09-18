package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/policy"
	"github.com/Marcosmxp/KryptaSec/internal/scan"
	"github.com/Marcosmxp/KryptaSec/internal/target"
)

type Service struct {
	Now func() time.Time
}

type StartScanRequest struct {
	Target     string
	ScopeHosts []string
}

type StartScanResult struct {
	Scan scan.Scan
}

func (s Service) StartScan(ctx context.Context, req StartScanRequest) (StartScanResult, error) {
	if err := ctx.Err(); err != nil {
		return StartScanResult{}, err
	}

	normalized, err := target.Normalize(req.Target)
	if err != nil {
		return StartScanResult{}, err
	}

	now := s.now()
	job, err := scan.New(string(normalized.Kind), normalized.Canonical, now)
	if err != nil {
		return StartScanResult{}, err
	}
	if err := job.Transition(scan.StatusValidatingScope, s.now()); err != nil {
		return StartScanResult{Scan: job}, err
	}

	if !(policy.Scope{Hosts: req.ScopeHosts}).Allows(normalized) {
		_ = job.Transition(scan.StatusFailed, s.now())
		return StartScanResult{Scan: job}, fmt.Errorf("target %q is outside the authorized scope", normalized.Canonical)
	}

	if err := job.Transition(scan.StatusReady, s.now()); err != nil {
		return StartScanResult{Scan: job}, err
	}
	return StartScanResult{Scan: job}, nil
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}
