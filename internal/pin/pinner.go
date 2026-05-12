// Package pin provides functionality to pin and retrieve specific
// environment variable sets by name and timestamp.
package pin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Pin represents a named snapshot of environment variables.
type Pin struct {
	Name      string            `json:"name"`
	Context   string            `json:"context"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages persisted pins on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating it if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("pin: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save persists a pin under its name, overwriting any existing pin.
func (s *Store) Save(p Pin) error {
	if p.Name == "" {
		return fmt.Errorf("pin: name must not be empty")
	}
	p.CreatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	return os.WriteFile(s.path(p.Name), data, 0o600)
}

// Load retrieves a pin by name.
func (s *Store) Load(name string) (Pin, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if os.IsNotExist(err) {
			return Pin{}, fmt.Errorf("pin %q not found", name)
		}
		return Pin{}, fmt.Errorf("pin: read: %w", err)
	}
	var p Pin
	if err := json.Unmarshal(data, &p); err != nil {
		return Pin{}, fmt.Errorf("pin: unmarshal: %w", err)
	}
	return p, nil
}

// Delete removes a pin by name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if os.IsNotExist(err) {
		return fmt.Errorf("pin %q not found", name)
	}
	return err
}

// List returns all stored pins sorted by creation time (oldest first).
func (s *Store) List() ([]Pin, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("pin: list dir: %w", err)
	}
	var pins []Pin
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		p, err := s.Load(e.Name())
		if err != nil {
			continue
		}
		pins = append(pins, p)
	}
	return pins, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}
