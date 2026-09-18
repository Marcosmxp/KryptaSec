package scan

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusCreated         Status = "created"
	StatusValidatingScope Status = "validating_scope"
	StatusReady           Status = "ready"
	StatusCancelled       Status = "cancelled"
	StatusFailed          Status = "failed"
)

type Scan struct {
	ID         string
	Target     string
	TargetKind string
	Status     Status
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func New(targetKind, target string, now time.Time) (Scan, error) {
	targetKind = strings.TrimSpace(targetKind)
	target = strings.TrimSpace(target)
	if targetKind == "" || target == "" {
		return Scan{}, fmt.Errorf("target kind and target are required")
	}

	id, err := newID()
	if err != nil {
		return Scan{}, fmt.Errorf("generate scan ID: %w", err)
	}
	now = now.UTC()
	return Scan{
		ID:         id,
		Target:     target,
		TargetKind: targetKind,
		Status:     StatusCreated,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (s *Scan) Transition(next Status, at time.Time) error {
	if s == nil {
		return fmt.Errorf("scan is nil")
	}
	if !allowedTransition(s.Status, next) {
		return fmt.Errorf("invalid scan transition %q -> %q", s.Status, next)
	}
	s.Status = next
	s.UpdatedAt = at.UTC()
	return nil
}

func allowedTransition(current, next Status) bool {
	switch current {
	case StatusCreated:
		return next == StatusValidatingScope || next == StatusCancelled || next == StatusFailed
	case StatusValidatingScope:
		return next == StatusReady || next == StatusCancelled || next == StatusFailed
	case StatusReady:
		return next == StatusCancelled
	default:
		return false
	}
}

func newID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "scan_" + hex.EncodeToString(b[:]), nil
}
