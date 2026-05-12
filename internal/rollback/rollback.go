package rollback

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a saved rollback point for a context.
type Entry struct {
	Context   string            `json:"context"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// Store manages rollback entries on disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("rollback: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save persists a rollback entry for the given context.
func (s *Store) Save(context string, vars map[string]string) error {
	if context == "" {
		return errors.New("rollback: context name must not be empty")
	}
	entry := Entry{
		Context:   context,
		Timestamp: time.Now().UTC(),
		Vars:      vars,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("rollback: marshal: %w", err)
	}
	filename := fmt.Sprintf("%s_%d.json", context, entry.Timestamp.UnixNano())
	return os.WriteFile(filepath.Join(s.dir, filename), data, 0600)
}

// List returns all rollback entries for the given context, sorted newest first.
func (s *Store) List(context string) ([]Entry, error) {
	pattern := filepath.Join(s.dir, context+"_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("rollback: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("rollback: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("rollback: unmarshal %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	return entries, nil
}

// Latest returns the most recent rollback entry for a context.
func (s *Store) Latest(context string) (*Entry, error) {
	entries, err := s.List(context)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("rollback: no entries found for context %q", context)
	}
	return &entries[0], nil
}
