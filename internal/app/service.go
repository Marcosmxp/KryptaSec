package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/policy"
	"github.com/Marcosmxp/KryptaSec/internal/scan"
	"github.com/Marcosmxp/KryptaSec/internal/store"
	"github.com/Marcosmxp/KryptaSec/internal/target"
)

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

type Service struct {
	Store  store.ScanRepository
	Logger *slog.Logger
	Now    func() time.Time
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
	s.log().DebugContext(
		ctx,
		"scan lifecycle transition",
		"scan_id", job.ID,
		"target_kind", job.TargetKind,
		"status", job.Status,
	)

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
	s.log().DebugContext(
		ctx,
		"scan lifecycle transition",
		"scan_id", job.ID,
		"target_kind", job.TargetKind,
		"status", job.Status,
	)

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
		s.log().WarnContext(
			ctx,
			"scan scope rejected",
			"scan_id", failed.ID,
			"target_kind", failed.TargetKind,
			"status", failed.Status,
		)
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
	s.log().DebugContext(
		ctx,
		"scan lifecycle transition",
		"scan_id", job.ID,
		"target_kind", job.TargetKind,
		"status", job.Status,
	)

	return StartScanResult{Scan: job}, nil
}

func (s Service) GetScan(ctx context.Context, id string) (scan.Scan, error) {
	if s.Store == nil {
		return scan.Scan{}, fmt.Errorf("scan store is required")
	}
	job, err := s.Store.Get(ctx, id)
	if err != nil {
		return scan.Scan{}, err
	}
	s.log().DebugContext(
		ctx,
		"scan loaded",
		"scan_id", job.ID,
		"target_kind", job.TargetKind,
		"status", job.Status,
	)
	return job, nil
}

func (s Service) log() *slog.Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return discardLogger
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}
