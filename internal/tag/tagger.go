package tag

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Tag associates a label with one or more context names.
type Tag struct {
	Name     string   `json:"name"`
	Contexts []string `json:"contexts"`
}

// Store manages tags persisted to a JSON file.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("tag: create dir: %w", err)
	}
	return &Store{path: path}, nil
}

func (s *Store) load() (map[string]Tag, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]Tag{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("tag: read: %w", err)
	}
	var tags map[string]Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("tag: parse: %w", err)
	}
	return tags, nil
}

func (s *Store) save(tags map[string]Tag) error {
	data, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return fmt.Errorf("tag: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Set creates or replaces a tag.
func (s *Store) Set(name string, contexts []string) error {
	if name == "" {
		return errors.New("tag: name must not be empty")
	}
	tags, err := s.load()
	if err != nil {
		return err
	}
	tags[name] = Tag{Name: name, Contexts: contexts}
	return s.save(tags)
}

// Get returns a tag by name.
func (s *Store) Get(name string) (Tag, error) {
	tags, err := s.load()
	if err != nil {
		return Tag{}, err
	}
	t, ok := tags[name]
	if !ok {
		return Tag{}, fmt.Errorf("tag: %q not found", name)
	}
	return t, nil
}

// List returns all tags sorted by name.
func (s *Store) List() ([]Tag, error) {
	tags, err := s.load()
	if err != nil {
		return nil, err
	}
	out := make([]Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Delete removes a tag by name.
func (s *Store) Delete(name string) error {
	tags, err := s.load()
	if err != nil {
		return err
	}
	if _, ok := tags[name]; !ok {
		return fmt.Errorf("tag: %q not found", name)
	}
	delete(tags, name)
	return s.save(tags)
}
