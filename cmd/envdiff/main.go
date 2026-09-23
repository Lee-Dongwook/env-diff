package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Lee-Dongwook/env-diff/internal/compare"
	"github.com/Lee-Dongwook/env-diff/internal/snapshot"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	switch args[0] {
	case "capture":
		return capture(args[1:])
	case "diff":
		return diffSnapshots(args[1:])
	default:
		printUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func capture(args []string) error {
	flags := flag.NewFlagSet("capture", flag.ContinueOnError)
	output := flags.String("out", "envdiff-snapshot.json", "output snapshot file")

	if err:= flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	result := snapshot.Capture()
	if err := snapshot.Write(*output, result); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}


	fmt.Printf("Environment snapshot saved to %s\n", *output)
	fmt.Printf("Recorded %d environment variables (values are not stored)\n", len(result.Environment))
	return nil
}

func diffSnapshots(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("diff requires two snapshot files")
	}

	left, err := snapshot.Read(args[0])
	if err != nil {
		return fmt.Errorf("read %s: %w", args[0], err)
	}

	right, err := snapshot.Read(args[1])
	if err != nil {
		return fmt.Errorf("read %s: %w", args[1], err)
	}

	result := compare.Diff(left, right)
	fmt.Printf("Comparing %s -> %s\n", args[0], args[1])
	hasDifferences := false

	if left.OS != right.OS {
		fmt.Printf("OS: %s -> %s\n", left.OS, right.OS)
		hasDifferences = true
	}
	if left.Architecture != right.Architecture {
		fmt.Printf("Architecture: %s -> %s\n", left.Architecture, right.Architecture)
		hasDifferences = true
	}
	if left.GoVersion != right.GoVersion {
		fmt.Printf("Go version: %s -> %s\n", left.GoVersion, right.GoVersion)
		hasDifferences = true
	}

	for _, name := range result.Added {
		fmt.Printf("+ Environment variable only in right: %s\n", name)
		hasDifferences = true
	}
	for _, name := range result.Removed {
		fmt.Printf("- Environment variable only in left: %s\n", name)
		hasDifferences = true
	}
	for _, name := range result.EmptyChanged {
		fmt.Printf("~ Environment variable empty status changed: %s\n", name)
		hasDifferences = true
	}

	if !hasDifferences {
		fmt.Println("No differences detected.")
	}

	return nil
}

func printUsage() {
	fmt.Println(`envdiff - compare local and CI environments

Usage:
  envdiff capture --out <file>
  envdiff diff <local-snapshot> <ci-snapshot>

Commands:
  capture  Capture the current environment
  diff     Compare two environment snapshots`)
}
