package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/Marcosmxp/KryptaSec/internal/app"
	"github.com/Marcosmxp/KryptaSec/internal/config"
	projectlogging "github.com/Marcosmxp/KryptaSec/internal/logging"
	"github.com/Marcosmxp/KryptaSec/internal/store"
	sqlitestore "github.com/Marcosmxp/KryptaSec/internal/store/sqlite"
)

type stringListFlag []string

func (s *stringListFlag) String() string {
	return fmt.Sprint([]string(*s))
}

func (s *stringListFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func runScan(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "status" {
		return runScanStatus(ctx, args[1:], stdout, stderr)
	}

	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var scopeHosts stringListFlag
	fs.Var(&scopeHosts, "scope-host", "exact remote hostname authorized for this scan; repeatable")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: kryptasec scan [--scope-host host] <target>")
		return 2
	}

	svc, closeStore, err := openScanService(stderr)
	if err != nil {
		fmt.Fprintf(stderr, "storage error: %v\n", err)
		return 1
	}
	defer closeStore()

	result, err := svc.StartScan(ctx, app.StartScanRequest{
		Target:     fs.Arg(0),
		ScopeHosts: scopeHosts,
	})
	if err != nil {
		fmt.Fprintf(stderr, "scan rejected: %v\n", err)
		return 1
	}

	writeScan(stdout, result.Scan.ID, result.Scan.Target, result.Scan.TargetKind, string(result.Scan.Status))
	fmt.Fprintln(stdout, "active testing: disabled (Phase 1)")
	return 0
}

func runScanStatus(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: kryptasec scan status <id>")
		return 2
	}

	svc, closeStore, err := openScanService(stderr)
	if err != nil {
		fmt.Fprintf(stderr, "storage error: %v\n", err)
		return 1
	}
	defer closeStore()

	job, err := svc.GetScan(ctx, args[0])
	if errors.Is(err, store.ErrNotFound) {
		fmt.Fprintf(stderr, "scan not found: %s\n", args[0])
		return 1
	}
	if err != nil {
		fmt.Fprintf(stderr, "scan status error: %v\n", err)
		return 1
	}

	writeScan(stdout, job.ID, job.Target, job.TargetKind, string(job.Status))
	fmt.Fprintf(stdout, "created: %s\n", job.CreatedAt.UTC().Format(time.RFC3339Nano))
	fmt.Fprintf(stdout, "updated: %s\n", job.UpdatedAt.UTC().Format(time.RFC3339Nano))
	return 0
}

func openScanService(logOutput io.Writer) (app.Service, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return app.Service{}, func() {}, err
	}

	s, err := sqlitestore.Open(filepath.Join(cfg.DataDir, "kryptasec.db"))
	if err != nil {
		return app.Service{}, func() {}, err
	}

	logger := projectlogging.New(logOutput, cfg.LogLevel)
	return app.Service{
		Store:  s,
		Logger: logger,
	}, func() { _ = s.Close() }, nil
}

func writeScan(w io.Writer, id, target, kind, status string) {
	fmt.Fprintf(w, "scan: %s\n", id)
	fmt.Fprintf(w, "target: %s\n", target)
	fmt.Fprintf(w, "kind: %s\n", kind)
	fmt.Fprintf(w, "status: %s\n", status)
}
