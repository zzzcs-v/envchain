package group

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// Group associates a named label with a set of context names.
type Group struct {
	Name     string   `json:"name"`
	Contexts []string `json:"contexts"`
}

// Store persists groups to a directory as JSON files.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating it if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Save persists g to disk, overwriting any existing entry.
func (s *Store) Save(g Group) error {
	if g.Name == "" {
		return errors.New("group name must not be empty")
	}
	if len(g.Contexts) == 0 {
		return errors.New("group must contain at least one context")
	}
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(g.Name), data, 0o640)
}

// Load retrieves the group with the given name.
func (s *Store) Load(name string) (Group, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Group{}, errors.New("group not found: " + name)
		}
		return Group{}, err
	}
	var g Group
	return g, json.Unmarshal(data, &g)
}

// Delete removes the group with the given name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("group not found: " + name)
	}
	return err
}

// List returns all saved groups sorted by name.
func (s *Store) List() ([]Group, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var groups []Group
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := e.Name()[:len(e.Name())-5]
		g, err := s.Load(name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups, nil
}
