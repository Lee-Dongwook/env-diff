package main

import (
	"fmt"
	"os"
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
		fmt.Println("capture command")
		return nil
	case "diff":
		fmt.Println("diff command")
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}
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
