package tag

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewStore(filepath.Join(dir, "tags.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSet_AndGet(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("infra", []string{"prod", "staging"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	tag, err := s.Get("infra")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if tag.Name != "infra" {
		t.Errorf("expected name infra, got %q", tag.Name)
	}
	if len(tag.Contexts) != 2 {
		t.Errorf("expected 2 contexts, got %d", len(tag.Contexts))
	}
}

func TestSet_EmptyName(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("", []string{"prod"}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestGet_Missing(t *testing.T) {
	s := tempStore(t)
	if _, err := s.Get("ghost"); err == nil {
		t.Error("expected error for missing tag")
	}
}

func TestList_SortedByName(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("zebra", []string{"dev"})
	_ = s.Set("alpha", []string{"prod"})
	_ = s.Set("beta", []string{"staging"})

	tags, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
	if tags[0].Name != "alpha" || tags[1].Name != "beta" || tags[2].Name != "zebra" {
		t.Errorf("unexpected order: %v", tags)
	}
}

func TestDelete_RemovesTag(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("temp", []string{"dev"})
	if err := s.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("temp"); err == nil {
		t.Error("expected tag to be deleted")
	}
}

func TestDelete_Missing(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("nope"); err == nil {
		t.Error("expected error deleting non-existent tag")
	}
}

func TestNewStore_InvalidPath(t *testing.T) {
	// Use a file as the directory component to force failure.
	f, _ := os.CreateTemp("", "envchain-tag-*")
	f.Close()
	defer os.Remove(f.Name())
	_, err := NewStore(filepath.Join(f.Name(), "tags.json"))
	if err == nil {
		t.Error("expected error for invalid store path")
	}
}
