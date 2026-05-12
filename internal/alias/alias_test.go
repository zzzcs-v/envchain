package alias

import (
	"os"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "alias-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSet_AndGet(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("prod", "production"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	a, err := s.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if a.Context != "production" {
		t.Errorf("expected context=production, got %q", a.Context)
	}
}

func TestSet_EmptyName(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("", "production"); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSet_EmptyContext(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("prod", ""); err == nil {
		t.Error("expected error for empty context")
	}
}

func TestGet_Missing(t *testing.T) {
	s := tempStore(t)
	if _, err := s.Get("nope"); err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestDelete_RemovesAlias(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("stg", "staging")
	if err := s.Delete("stg"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("stg"); err == nil {
		t.Error("expected alias to be gone")
	}
}

func TestDelete_Missing(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error deleting missing alias")
	}
}

func TestList_SortedByName(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("zz", "ctx-z")
	_ = s.Set("aa", "ctx-a")
	_ = s.Set("mm", "ctx-m")
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(list))
	}
	if list[0].Name != "aa" || list[1].Name != "mm" || list[2].Name != "zz" {
		t.Errorf("wrong sort order: %v", list)
	}
}
