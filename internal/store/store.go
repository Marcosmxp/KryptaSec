package store

import (
	"context"
	"errors"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/scan"
)

var (
	ErrNotFound = errors.New("scan not found")
	ErrConflict = errors.New("scan state conflict")
)

type ScanRepository interface {
	Create(ctx context.Context, job scan.Scan) error
	Get(ctx context.Context, id string) (scan.Scan, error)
	Transition(ctx context.Context, id string, expected scan.Status, next scan.Status, at time.Time) (scan.Scan, error)
}
