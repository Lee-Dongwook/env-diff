package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	excludePrefix := flags.String("exclude-prefix", "", "comma-separated environment variable prefixes to exclude",)

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	result := snapshot.Capture(parsePrefixes(*excludePrefix)...)
	if err := snapshot.Write(*output, result); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}

	fmt.Printf("Environment snapshot saved to %s\n", *output)
	fmt.Printf("Recorded %d environment variables and %d toolchains (values are not stored)\n",
		len(result.Environment),
		len(result.Toolchains),
	)
	return nil
}

func diffSnapshots(args []string) error {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	failOnDiff := flags.Bool("fail-on-diff", false, "exit with an error when differences are found")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) != 2 {
		return fmt.Errorf("diff requires two snapshot files")
	}

	left, err := snapshot.Read(files[0])
	if err != nil {
		return fmt.Errorf("read %s: %w", files[0], err)
	}

	right, err := snapshot.Read(files[1])
	if err != nil {
		return fmt.Errorf("read %s: %w",files[1], err)
	}

	result := compare.Diff(left, right)
	fmt.Printf("Comparing %s -> %s\n", files[0], files[1])

	hasDifferences := false

	if left.OS != right.OS {
		fmt.Printf("OS: %s -> %s\n", left.OS, right.OS)
		hasDifferences = true
	}
	if left.Architecture != right.Architecture {
		fmt.Printf("Architecture: %s -> %s\n", left.Architecture, right.Architecture)
		hasDifferences = true
	}

	for _, change := range result.ToolchainChanges {
		fmt.Printf("%s: %s -> %s\n", change.Name, change.Left, change.Right)
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

	if hasDifferences && *failOnDiff {
		return fmt.Errorf("differences detected")
	}

	return nil
}

func parsePrefixes(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	var prefixes []string
	for _, prefix := range strings.Split(value,  ",") {
		prefix = strings.TrimSpace(prefix)
		if prefix != "" {
			prefixes = append(prefixes, prefix)
		}
	}

	return prefixes
}

func printUsage() {
	fmt.Println(`envdiff - compare local and CI environments

Usage:
  envdiff capture --out <file>
  envdiff diff <local-snapshot> <ci-snapshot>
  envdiff diff [--fail-on-diff] <local-snapshot> <ci-snapshot>

Commands:
  capture  Capture the current environment
  diff     Compare two environment snapshots`)
}
