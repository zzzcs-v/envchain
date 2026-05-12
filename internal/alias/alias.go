// Package alias manages short-name aliases for environment contexts.
package alias

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Alias maps a short name to a context name.
type Alias struct {
	Name    string `json:"name"`
	Context string `json:"context"`
}

// Store persists aliases to disk.
type Store struct {
	dir string
}

// NewStore creates a new Store backed by dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("alias: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Set stores an alias.
func (s *Store) Set(name, context string) error {
	if name == "" {
		return errors.New("alias: name must not be empty")
	}
	if context == "" {
		return errors.New("alias: context must not be empty")
	}
	a := Alias{Name: name, Context: context}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("alias: marshal: %w", err)
	}
	return os.WriteFile(s.path(name), data, 0o600)
}

// Get retrieves an alias by name.
func (s *Store) Get(name string) (*Alias, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("alias: %q not found", name)
		}
		return nil, fmt.Errorf("alias: read: %w", err)
	}
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("alias: unmarshal: %w", err)
	}
	return &a, nil
}

// Delete removes an alias.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("alias: %q not found", name)
	}
	return err
}

// List returns all aliases sorted by name.
func (s *Store) List() ([]Alias, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("alias: list: %w", err)
	}
	var aliases []Alias
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := e.Name()[:len(e.Name())-5]
		a, err := s.Get(name)
		if err != nil {
			continue
		}
		aliases = append(aliases, *a)
	}
	sort.Slice(aliases, func(i, j int) bool { return aliases[i].Name < aliases[j].Name })
	return aliases, nil
}
