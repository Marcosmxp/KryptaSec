package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/scan"
	"github.com/Marcosmxp/KryptaSec/internal/store"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	s := &Store{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	if err := s.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Create(ctx context.Context, job scan.Scan) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO scans (id, target, target_kind, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		job.ID,
		job.Target,
		job.TargetKind,
		string(job.Status),
		job.CreatedAt.UTC().Format(time.RFC3339Nano),
		job.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("create scan: %w", err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (scan.Scan, error) {
	return getScan(ctx, s.db, id)
}

func (s *Store) Transition(
	ctx context.Context,
	id string,
	expected scan.Status,
	next scan.Status,
	at time.Time,
) (scan.Scan, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return scan.Scan{}, fmt.Errorf("begin scan transition: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	current, err := getScan(ctx, tx, id)
	if err != nil {
		return scan.Scan{}, err
	}
	if current.Status != expected {
		return scan.Scan{}, fmt.Errorf(
			"%w: scan %s is %q, expected %q",
			store.ErrConflict,
			id,
			current.Status,
			expected,
		)
	}

	if err := current.Transition(next, at); err != nil {
		return scan.Scan{}, err
	}

	result, err := tx.ExecContext(
		ctx,
		"UPDATE scans SET status = ?, updated_at = ? WHERE id = ? AND status = ?",
		string(current.Status),
		current.UpdatedAt.UTC().Format(time.RFC3339Nano),
		id,
		string(expected),
	)
	if err != nil {
		return scan.Scan{}, fmt.Errorf("update scan status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return scan.Scan{}, fmt.Errorf("read transition result: %w", err)
	}
	if affected != 1 {
		return scan.Scan{}, fmt.Errorf("%w: scan %s changed concurrently", store.ErrConflict, id)
	}

	if err := tx.Commit(); err != nil {
		return scan.Scan{}, fmt.Errorf("commit scan transition: %w", err)
	}
	return current, nil
}

func (s *Store) migrate(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin sqlite migration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS scans (
			id TEXT PRIMARY KEY,
			target TEXT NOT NULL,
			target_kind TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_scans_status ON scans(status);
	`); err != nil {
		return fmt.Errorf("apply sqlite migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit sqlite migration: %w", err)
	}
	return nil
}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func getScan(ctx context.Context, q queryer, id string) (scan.Scan, error) {
	var (
		job       scan.Scan
		status    string
		createdAt string
		updatedAt string
	)

	err := q.QueryRowContext(
		ctx,
		"SELECT id, target, target_kind, status, created_at, updated_at FROM scans WHERE id = ?",
		id,
	).Scan(
		&job.ID,
		&job.Target,
		&job.TargetKind,
		&status,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return scan.Scan{}, fmt.Errorf("%w: %s", store.ErrNotFound, id)
	}
	if err != nil {
		return scan.Scan{}, fmt.Errorf("get scan: %w", err)
	}

	job.Status = scan.Status(status)
	job.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return scan.Scan{}, fmt.Errorf("parse scan created_at: %w", err)
	}
	job.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return scan.Scan{}, fmt.Errorf("parse scan updated_at: %w", err)
	}

	return job, nil
}
