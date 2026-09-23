package compare

import "github.com/Lee-Dongwook/env-diff/internal/snapshot"

type Result struct {
	OSChanged           bool
	ArchitectureChanged bool
	GoVersionChanged    bool
	Added               []string
	Removed             []string
	EmptyChanged        []string
}

func Diff(left, right snapshot.Snapshot) Result {
	result := Result{
		OSChanged:           left.OS != right.OS,
		ArchitectureChanged: left.Architecture != right.Architecture,
		GoVersionChanged:    left.GoVersion != right.GoVersion,
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

	return result
}
