package namespace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Entry represents a namespace binding: a named group of context keys.
type Entry struct {
	Name     string            `json:"name"`
	Prefix   string            `json:"prefix"`
	Contexts []string          `json:"contexts"`
	Meta     map[string]string `json:"meta,omitempty"`
}

// Store manages namespace entries on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Save persists an Entry to disk.
func (s *Store) Save(e Entry) error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("namespace name must not be empty")
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(e.Name), data, 0o600)
}

// Load retrieves an Entry by name.
func (s *Store) Load(name string) (Entry, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, errors.New("namespace not found: " + name)
		}
		return Entry{}, err
	}
	var e Entry
	return e, json.Unmarshal(data, &e)
}

// Delete removes an Entry by name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("namespace not found: " + name)
	}
	return err
}

// List returns all stored entries sorted by name.
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var out []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, err
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
