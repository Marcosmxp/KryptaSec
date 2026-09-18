package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/Marcosmxp/KryptaSec/internal/app"
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

	result, err := (app.Service{}).StartScan(ctx, app.StartScanRequest{
		Target:     fs.Arg(0),
		ScopeHosts: scopeHosts,
	})
	if err != nil {
		fmt.Fprintf(stderr, "scan rejected: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "scan: %s\n", result.Scan.ID)
	fmt.Fprintf(stdout, "target: %s\n", result.Scan.Target)
	fmt.Fprintf(stdout, "kind: %s\n", result.Scan.TargetKind)
	fmt.Fprintf(stdout, "status: %s\n", result.Scan.Status)
	fmt.Fprintln(stdout, "active testing: disabled (Phase 1)")
	return 0
}
