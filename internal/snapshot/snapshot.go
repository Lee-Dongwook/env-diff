package snapshot

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Snapshot struct {
	SchemaVersion int           `json:"schemaVersion"`
	CapturedAt    time.Time     `json:"capturedAt"`
	OS            string        `json:"os"`
	Architecture  string        `json:"architecture"`
	Toolchains    []ToolVersion `json:"toolchains"`
	Environment   []Environment `json:"environment"`
}

type ToolVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Environment struct {
	Name  string `json:"name"`
	Empty bool   `json:"empty"`
}

func Capture() Snapshot {
	values := make(map[string]string)

	for _, entry := range os.Environ() {
		name, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		if shouldIgnoreEnvironmentVariable(name) {
			continue
		}

		values[name] = value
	}

	environment := make([]Environment, 0, len(values))
	for name, value := range values {
		environment = append(environment, Environment{
			Name:  name,
			Empty: value == "",
		})
	}

	sort.Slice(environment, func(i, j int) bool {
		return environment[i].Name < environment[j].Name
	})

	toolchains := captureToolchains()

	return Snapshot{
		SchemaVersion: 1,
		CapturedAt:    time.Now().UTC(),
		OS:            runtime.GOOS,
		Architecture:  runtime.GOARCH,
		Toolchains:    toolchains,
		Environment:   environment,
	}
}

func Write(path string, snapshot Snapshot) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(snapshot)
}

func Read(path string) (Snapshot, error) {
	file, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer file.Close()

	var result Snapshot
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return Snapshot{}, err
	}

	return result, nil
}

func captureToolchains() []ToolVersion {
	probes := []struct {
		name string
		args []string
	}{
		{name: "node", args: []string{"--version"}},
		{name: "npm", args: []string{"--version"}},
		{name: "go", args: []string{"version"}},
		{name: "python3", args: []string{"--version"}},
	}

	var versions []ToolVersion

	for _, probe := range probes {
		path, err := exec.LookPath(probe.name)
		if err != nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		output, err := exec.CommandContext(ctx, path, probe.args...).Output()
		cancel()
		if err != nil {
			continue
		}

		versions = append(versions, ToolVersion{
			Name:    probe.name,
			Version: strings.TrimSpace(string(output)),
		})
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Name < versions[j].Name
	})

	return versions
}

func shouldIgnoreEnvironmentVariable(name string) bool {
	upperName := strings.ToUpper(name)

	if upperName == "CI" {
		return true
	}

	ignoredPrefixes := []string{
		"GITHUB_",
		"ACTIONS_",
		"RUNNER_",
	}

	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(upperName, prefix) {
			return true
		}
	}

	return false
}
