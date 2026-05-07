package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a saved snapshot of resolved env vars for a context.
type Entry struct {
	Context   string            `json:"context"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages snapshot files on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that persists snapshots under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot for the given context and vars.
func (s *Store) Save(context string, vars map[string]string) error {
	entry := Entry{
		Context:   context,
		Vars:      vars,
		CreatedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := s.filePath(context)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load reads the latest snapshot for the given context.
func (s *Store) Load(context string) (*Entry, error) {
	path := s.filePath(context)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: no snapshot found for context %q", context)
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &entry, nil
}

// Delete removes the snapshot for the given context.
func (s *Store) Delete(context string) error {
	path := s.filePath(context)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}

func (s *Store) filePath(context string) string {
	return filepath.Join(s.dir, context+".json")
}
