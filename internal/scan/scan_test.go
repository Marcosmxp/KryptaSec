package scan

import (
	"testing"
	"time"
)

func TestNewStartsCreated(t *testing.T) {
	now := time.Date(2026, 9, 18, 20, 0, 0, 0, time.UTC)
	s, err := New("local", "/tmp/app", now)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if s.ID == "" {
		t.Fatal("expected non-empty scan ID")
	}
	if s.Status != StatusCreated {
		t.Fatalf("Status = %q, want %q", s.Status, StatusCreated)
	}
	if !s.CreatedAt.Equal(now) || !s.UpdatedAt.Equal(now) {
		t.Fatalf("unexpected timestamps: created=%v updated=%v", s.CreatedAt, s.UpdatedAt)
	}
}

func TestTransitionAllowsPhaseOneLifecycle(t *testing.T) {
	now := time.Now().UTC()
	s, err := New("http", "https://example.com", now)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Transition(StatusValidatingScope, now.Add(time.Second)); err != nil {
		t.Fatalf("created -> validating_scope: %v", err)
	}
	if err := s.Transition(StatusReady, now.Add(2*time.Second)); err != nil {
		t.Fatalf("validating_scope -> ready: %v", err)
	}
	if s.Status != StatusReady {
		t.Fatalf("Status = %q, want ready", s.Status)
	}
}

func TestTransitionRejectsInvalidStateChange(t *testing.T) {
	now := time.Now().UTC()
	s, err := New("local", "/tmp/app", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Transition(StatusReady, now.Add(time.Second)); err == nil {
		t.Fatal("created -> ready should be rejected")
	}
}

func TestReadyCanBeCancelled(t *testing.T) {
	now := time.Now().UTC()
	s, _ := New("local", "/tmp/app", now)
	_ = s.Transition(StatusValidatingScope, now.Add(time.Second))
	_ = s.Transition(StatusReady, now.Add(2*time.Second))

	if err := s.Transition(StatusCancelled, now.Add(3*time.Second)); err != nil {
		t.Fatalf("ready -> cancelled: %v", err)
	}
}
