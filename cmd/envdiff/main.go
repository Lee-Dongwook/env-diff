package main

import (
	"flag"
	"fmt"
	"os"

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
		fmt.Println("diff command")
		return nil
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

func printUsage() {
	fmt.Println(`envdiff - compare local and CI environments

Usage:
  envdiff capture --out <file>
  envdiff diff <local-snapshot> <ci-snapshot>

Commands:
  capture  Capture the current environment
  diff     Compare two environment snapshots`)
}
