package compare

import (
	"sort"

	"github.com/Lee-Dongwook/env-diff/internal/snapshot"
)

type ToolchainChange struct {
	Name  string
	Left  string
	Right string
}

type Result struct {
	OSChanged           bool
	ArchitectureChanged bool
	GoVersionChanged    bool
	Added               []string
	Removed             []string
	EmptyChanged        []string
	ToolchainChanges    []ToolchainChange
}

func Diff(left, right snapshot.Snapshot) Result {
	result := Result{
		OSChanged:           left.OS != right.OS,
		ArchitectureChanged: left.Architecture != right.Architecture,
		ToolchainChanges:    []ToolchainChange{},
		Added:               []string{},
		Removed:             []string{},
		EmptyChanged:        []string{},
	}

	leftEnvironment := make(map[string]bool, len(left.Environment))
	for _, variable := range left.Environment {
		leftEnvironment[variable.Name] = variable.Empty
	}

	rightEnvironment := make(map[string]bool, len(right.Environment))
	for _, variable := range right.Environment {
		rightEnvironment[variable.Name] = variable.Empty
	}

	for name, leftEmpty := range leftEnvironment {
		rightEmpty, exists := rightEnvironment[name]

		if !exists {
			result.Removed = append(result.Removed, name)
			continue
		}

		if leftEmpty != rightEmpty {
			result.EmptyChanged = append(result.EmptyChanged, name)
		}
	}

	for name := range rightEnvironment {
		if _, exists := leftEnvironment[name]; !exists {
			result.Added = append(result.Added, name)
		}
	}

	leftTools := make(map[string]string, len(left.Toolchains))
	for _, tool := range left.Toolchains {
		leftTools[tool.Name] = tool.Version
	}

	rightTools := make(map[string]string, len(right.Toolchains))
	for _, tool := range right.Toolchains {
		rightTools[tool.Name] = tool.Version
	}

	toolNames := make(map[string]struct{})
	for name := range leftTools {
		toolNames[name] = struct{}{}
	}
	for name := range rightTools {
		toolNames[name] = struct{}{}
	}

	for name := range toolNames {
		leftVersion, leftExists := leftTools[name]
		rightVersion, rightExists := rightTools[name]

		if !leftExists {
			leftVersion = "not installed"
		}
		if !rightExists {
			rightVersion = "not installed"
		}

		if leftVersion != rightVersion {
			result.ToolchainChanges = append(result.ToolchainChanges, ToolchainChange{
				Name:  name,
				Left:  leftVersion,
				Right: rightVersion,
			})
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.EmptyChanged)
	sort.Slice(result.ToolchainChanges, func(i, j int) bool {
		return result.ToolchainChanges[i].Name < result.ToolchainChanges[j].Name
	})

	return result
}
