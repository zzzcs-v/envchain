package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Profile represents a named set of context + export preferences.
type Profile struct {
	Name    string            `json:"name"`
	Context string            `json:"context"`
	Format  string            `json:"format"`
	Extras  map[string]string `json:"extras,omitempty"`
}

// Store manages profiles on disk.
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

// Save writes a profile to disk.
func (s *Store) Save(p Profile) error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(p.Name), data, 0o600)
}

// Load reads a profile by name from disk.
func (s *Store) Load(name string) (Profile, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Profile{}, errors.New("profile not found: " + name)
		}
		return Profile{}, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

// Delete removes a profile from disk.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("profile not found: " + name)
	}
	return err
}

// List returns all profile names stored on disk.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Exists reports whether a profile with the given name exists on disk.
func (s *Store) Exists(name string) (bool, error) {
	_, err := os.Stat(s.path(name))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
