package snapshot

import (
	"encoding/json"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Snapshot struct {
	SchemaVersion int  			`json:"schemaVersion"`
	CapturedAt 	  time.Time		`json:"capturedAt"`
	OS			  string		`json:"os"`
	Architecture  string		`json:"architecture"`
	GoVersion	  string		`json:"goVersion"`
	Environment	  []Environment	`json:"environment"`
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
		values[name] = value
	}

	environment := make([]Environment, 0, len(values))
	for name, value := range values {
		environment = append(environment, Environment{
			Name: name,
			Empty: value == "",
		})
	}

	sort.Slice(environment, func(i, j int) bool {
		return environment[i].Name < environment[j].Name
	})

	return Snapshot{
		SchemaVersion: 1,
		CapturedAt: time.Now().UTC(),
		OS: runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion: runtime.Version(),
		Environment: environment,
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
