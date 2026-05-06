package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for environment variables.
type Format string

const (
	FormatDotenv  Format = "dotenv"
	FormatExport  Format = "export"
	FormatJSON    Format = "json"
)

// Exporter writes resolved env vars to a writer in a given format.
type Exporter struct {
	format Format
}

// NewExporter creates a new Exporter for the given format.
func NewExporter(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatExport, FormatJSON:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q: must be one of dotenv, export, json", format)
	}
}

// Write serializes the env map to the writer in the configured format.
func (e *Exporter) Write(w io.Writer, env map[string]string) error {
	keys := sortedKeys(env)
	switch e.format {
	case FormatDotenv:
		return writeDotenv(w, keys, env)
	case FormatExport:
		return writeExport(w, keys, env)
	case FormatJSON:
		return writeJSON(w, keys, env)
	}
	return nil
}

func writeDotenv(w io.Writer, keys []string, env map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, quoteValue(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, keys []string, env map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, quoteValue(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, keys []string, env map[string]string) error {
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("  %q: %q", k, env[k]))
	}
	_, err := fmt.Fprintf(w, "{\n%s\n}\n", strings.Join(lines, ",\n"))
	return err
}

func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
