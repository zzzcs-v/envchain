package chain

import (
	"testing"
)

func TestSet_AndGet(t *testing.T) {
	s := New()
	err := s.Set("mychain", []string{"dev", "staging", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, err := s.Get("mychain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Contexts) != 3 || c.Contexts[1] != "staging" {
		t.Errorf("unexpected contexts: %v", c.Contexts)
	}
}

func TestSet_EmptyName(t *testing.T) {
	s := New()
	if err := s.Set("", []string{"dev"}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSet_NoContexts(t *testing.T) {
	s := New()
	if err := s.Set("empty", []string{}); err == nil {
		t.Error("expected error for empty contexts")
	}
}

func TestSet_DuplicateContext(t *testing.T) {
	s := New()
	if err := s.Set("dup", []string{"dev", "dev"}); err == nil {
		t.Error("expected error for duplicate context")
	}
}

func TestSet_EmptyContextName(t *testing.T) {
	s := New()
	if err := s.Set("bad", []string{"dev", ""}); err == nil {
		t.Error("expected error for empty context name in list")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	if _, err := s.Get("nope"); err == nil {
		t.Error("expected error for missing chain")
	}
}

func TestDelete_RemovesChain(t *testing.T) {
	s := New()
	_ = s.Set("c", []string{"dev"})
	if err := s.Delete("c"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Get("c"); err == nil {
		t.Error("expected chain to be deleted")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error for missing chain")
	}
}

func TestList_SortedByName(t *testing.T) {
	s := New()
	_ = s.Set("zebra", []string{"prod"})
	_ = s.Set("alpha", []string{"dev"})
	_ = s.Set("middle", []string{"staging"})
	list := s.List()
	if len(list) != 3 {
		t.Fatalf("expected 3, got %d", len(list))
	}
	if list[0].Name != "alpha" || list[1].Name != "middle" || list[2].Name != "zebra" {
		t.Errorf("unexpected order: %v", list)
	}
}
