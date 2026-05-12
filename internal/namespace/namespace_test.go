package namespace

import (
	"os"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "ns-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSave_AndLoad(t *testing.T) {
	s := tempStore(t)
	e := Entry{Name: "infra", Prefix: "INFRA_", Contexts: []string{"dev", "prod"}}
	if err := s.Save(e); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load("infra")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != e.Name || got.Prefix != e.Prefix {
		t.Errorf("got %+v, want %+v", got, e)
	}
}

func TestSave_EmptyName(t *testing.T) {
	s := tempStore(t)
	if err := s.Save(Entry{Name: ""}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLoad_Missing(t *testing.T) {
	s := tempStore(t)
	_, err := s.Load("ghost")
	if err == nil {
		t.Error("expected error for missing namespace")
	}
}

func TestDelete_RemovesNamespace(t *testing.T) {
	s := tempStore(t)
	e := Entry{Name: "temp", Prefix: "T_", Contexts: []string{"dev"}}
	_ = s.Save(e)
	if err := s.Delete("temp"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete("temp"); err == nil {
		t.Error("expected error deleting already-removed namespace")
	}
}

func TestList_SortedByName(t *testing.T) {
	s := tempStore(t)
	for _, name := range []string{"zebra", "alpha", "mango"} {
		_ = s.Save(Entry{Name: name, Prefix: name + "_", Contexts: []string{"dev"}})
	}
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3, got %d", len(list))
	}
	if list[0].Name != "alpha" || list[1].Name != "mango" || list[2].Name != "zebra" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestList_Empty(t *testing.T) {
	s := tempStore(t)
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}
