package profile

import (
	"os"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSave_AndLoad(t *testing.T) {
	s := tempStore(t)
	p := Profile{Name: "dev", Context: "development", Format: "dotenv"}
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Context != "development" || got.Format != "dotenv" {
		t.Errorf("unexpected profile: %+v", got)
	}
}

func TestSave_EmptyName(t *testing.T) {
	s := tempStore(t)
	err := s.Save(Profile{Name: "", Context: "dev"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestLoad_Missing(t *testing.T) {
	s := tempStore(t)
	_, err := s.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	s := tempStore(t)
	_ = s.Save(Profile{Name: "staging", Context: "staging", Format: "export"})
	if err := s.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Load("staging")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDelete_Missing(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("nope"); err == nil {
		t.Fatal("expected error deleting missing profile")
	}
}

func TestList_ReturnsNames(t *testing.T) {
	s := tempStore(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = s.Save(Profile{Name: name, Context: name, Format: "json"})
	}
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(names))
	}
}

func TestList_EmptyDir(t *testing.T) {
	s := tempStore(t)
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/profiles"
	_, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
}
