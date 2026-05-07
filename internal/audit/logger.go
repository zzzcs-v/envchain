package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Action    string            `json:"action"`
	Context   string            `json:"context"`
	Vars      []string          `json:"vars,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes audit entries to a destination.
type Logger struct {
	path string
	f    *os.File
}

// NewLogger opens (or creates) the audit log file at path.
func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{path: path, f: f}, nil
}

// Log writes a single audit entry as a JSON line.
func (l *Logger) Log(action, context string, vars []string, meta map[string]string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Context:   context,
		Vars:      vars,
		Meta:      meta,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.f, "%s\n", data)
	return err
}

// Close closes the underlying log file.
func (l *Logger) Close() error {
	return l.f.Close()
}

// ReadAll parses all entries from the log file at path.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
