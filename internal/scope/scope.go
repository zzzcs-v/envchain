package scope

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"encoding/json"
)

// Scope represents a named boundary grouping environment variables under a label.
type Scope struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Store manages scopes on disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("scope: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Save persists a scope to disk.
func (s *Store) Save(sc Scope) error {
	if sc.Name == "" {
		return errors.New("scope: name must not be empty")
	}
	data, err := json.MarshalIndent(sc, "", "  ")
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	return os.WriteFile(s.path(sc.Name), data, 0600)
}

// Load retrieves a scope by name.
func (s *Store) Load(name string) (Scope, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Scope{}, fmt.Errorf("scope %q not found", name)
		}
		return Scope{}, fmt.Errorf("scope: read: %w", err)
	}
	var sc Scope
	if err := json.Unmarshal(data, &sc); err != nil {
		return Scope{}, fmt.Errorf("scope: unmarshal: %w", err)
	}
	return sc, nil
}

// Delete removes a scope by name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("scope %q not found", name)
	}
	return err
}

// List returns all scope names sorted alphabetically.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("scope: list: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	sort.Strings(names)
	return names, nil
}
