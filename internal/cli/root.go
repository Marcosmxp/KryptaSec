package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/Marcosmxp/KryptaSec/internal/config"
	"github.com/Marcosmxp/KryptaSec/internal/doctor"
)

func Run(ctx context.Context, args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 1 && (args[0] == "--version" || args[0] == "version") {
		fmt.Fprintf(stdout, "KryptaSec %s\n", version)
		return 0
	}

	if len(args) == 0 {
		writeUsage(stdout)
		return 0
	}

	if args[0] == "doctor" {
		return runDoctor(ctx, stdout, stderr)
	}
	if args[0] == "scan" {
		return runScan(ctx, args[1:], stdout, stderr)
	}

	fmt.Fprintf(stderr, "unknown command %q\n", args[0])
	writeUsage(stderr)
	return 2
}

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: kryptasec <command> [options]")
	fmt.Fprintln(w, "Commands: doctor, scan, version")
}

func runDoctor(ctx context.Context, stdout, stderr io.Writer) int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(stderr, "configuration error: %v\n", err)
		return 1
	}

	report := doctor.Run(ctx, cfg)
	for _, check := range report.Checks {
		status := "ok"
		if !check.OK {
			status = "fail"
		}
		fmt.Fprintf(stdout, "[%s] %s: %s\n", status, check.Name, check.Message)
	}
	if !report.OK() {
		return 1
	}
	return 0
}
