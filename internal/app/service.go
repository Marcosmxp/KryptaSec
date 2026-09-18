package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/policy"
	"github.com/Marcosmxp/KryptaSec/internal/scan"
	"github.com/Marcosmxp/KryptaSec/internal/store"
	"github.com/Marcosmxp/KryptaSec/internal/target"
)

type Service struct {
	Store store.ScanRepository
	Now   func() time.Time
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
	if s.Store == nil {
		return StartScanResult{}, fmt.Errorf("scan store is required")
	}

	normalized, err := target.Normalize(req.Target)
	if err != nil {
		return StartScanResult{}, err
	}

	job, err := scan.New(string(normalized.Kind), normalized.Canonical, s.now())
	if err != nil {
		return StartScanResult{}, err
	}
	if err := s.Store.Create(ctx, job); err != nil {
		return StartScanResult{}, err
	}

	job, err = s.Store.Transition(
		ctx,
		job.ID,
		scan.StatusCreated,
		scan.StatusValidatingScope,
		s.now(),
	)
	if err != nil {
		return StartScanResult{Scan: job}, err
	}

	if !(policy.Scope{Hosts: req.ScopeHosts}).Allows(normalized) {
		failed, transitionErr := s.Store.Transition(
			ctx,
			job.ID,
			scan.StatusValidatingScope,
			scan.StatusFailed,
			s.now(),
		)
		if transitionErr != nil {
			return StartScanResult{Scan: job}, transitionErr
		}
		return StartScanResult{Scan: failed}, fmt.Errorf(
			"target %q is outside the authorized scope",
			normalized.Canonical,
		)
	}

	job, err = s.Store.Transition(
		ctx,
		job.ID,
		scan.StatusValidatingScope,
		scan.StatusReady,
		s.now(),
	)
	if err != nil {
		return StartScanResult{Scan: job}, err
	}

	return StartScanResult{Scan: job}, nil
}

func (s Service) GetScan(ctx context.Context, id string) (scan.Scan, error) {
	if s.Store == nil {
		return scan.Scan{}, fmt.Errorf("scan store is required")
	}
	return s.Store.Get(ctx, id)
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}
