package group_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain/internal/group"
)

func tempStore(t *testing.T) *group.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "group-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := group.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSave_AndLoad(t *testing.T) {
	s := tempStore(t)
	g := group.Group{Name: "backend", Contexts: []string{"dev", "staging"}}
	if err := s.Save(g); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("backend")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != g.Name || len(got.Contexts) != 2 {
		t.Errorf("unexpected group: %+v", got)
	}
}

func TestSave_EmptyName(t *testing.T) {
	s := tempStore(t)
	if err := s.Save(group.Group{Contexts: []string{"dev"}}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSave_EmptyContexts(t *testing.T) {
	s := tempStore(t)
	if err := s.Save(group.Group{Name: "empty"}); err == nil {
		t.Error("expected error for empty contexts")
	}
}

func TestLoad_Missing(t *testing.T) {
	s := tempStore(t)
	if _, err := s.Load("ghost"); err == nil {
		t.Error("expected error for missing group")
	}
}

func TestDelete_RemovesGroup(t *testing.T) {
	s := tempStore(t)
	g := group.Group{Name: "infra", Contexts: []string{"prod"}}
	_ = s.Save(g)
	if err := s.Delete("infra"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Load("infra"); err == nil {
		t.Error("expected group to be gone")
	}
}

func TestList_SortedByName(t *testing.T) {
	s := tempStore(t)
	for _, name := range []string{"zebra", "alpha", "mango"} {
		_ = s.Save(group.Group{Name: name, Contexts: []string{"dev"}})
	}
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(list))
	}
	if list[0].Name != "alpha" || list[1].Name != "mango" || list[2].Name != "zebra" {
		t.Errorf("unexpected order: %v", list)
	}
}
