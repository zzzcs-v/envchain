package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single history record for a context resolution.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Context   string            `json:"context"`
	Format    string            `json:"format"`
	Vars      map[string]string `json:"vars"`
}

// Store manages history entries on disk.
type Store struct {
	dir string
}

// NewStore creates a new Store backed by the given directory.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("history: mkdir %s: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Record saves an entry to disk.
func (s *Store) Record(ctx, format string, vars map[string]string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Context:   ctx,
		Format:    format,
		Vars:      vars,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	name := fmt.Sprintf("%d.json", e.Timestamp.UnixNano())
	path := filepath.Join(s.dir, name)
	return os.WriteFile(path, data, 0o600)
}

// List returns all entries sorted by timestamp ascending.
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

// Clear removes all history entries.
func (s *Store) Clear() error {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("history: glob: %w", err)
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil {
			return fmt.Errorf("history: remove %s: %w", m, err)
		}
	}
	return nil
}
