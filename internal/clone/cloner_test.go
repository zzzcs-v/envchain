package clone

import (
	"testing"
)

func makeSource() map[string][]Entry {
	return map[string][]Entry{
		"dev": {
			{Key: "DB_HOST", Value: "localhost"},
			{Key: "DB_PORT", Value: "5432"},
			{Key: "API_KEY", Value: "dev-secret"},
		},
		"staging": {
			{Key: "DB_HOST", Value: "staging.db"},
		},
	}
}

func TestNew_NilSource(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil source, got nil")
	}
}

func TestClone_FullCopy(t *testing.T) {
	c, _ := New(makeSource())
	ctx, err := c.Clone("dev", "dev-copy", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Name != "dev-copy" {
		t.Errorf("expected name dev-copy, got %s", ctx.Name)
	}
	if len(ctx.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ctx.Entries))
	}
}

func TestClone_WithFilter(t *testing.T) {
	c, _ := New(makeSource())
	ctx, err := c.Clone("dev", "dev-filtered", []string{"DB_HOST", "DB_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(ctx.Entries))
	}
	for _, e := range ctx.Entries {
		if e.Key == "API_KEY" {
			t.Error("API_KEY should have been filtered out")
		}
	}
}

func TestClone_MissingSource(t *testing.T) {
	c, _ := New(makeSource())
	_, err := c.Clone("prod", "prod-copy", nil)
	if err == nil {
		t.Fatal("expected error for missing source context")
	}
}

func TestClone_EmptyDestination(t *testing.T) {
	c, _ := New(makeSource())
	_, err := c.Clone("dev", "", nil)
	if err == nil {
		t.Fatal("expected error for empty destination name")
	}
}

func TestClone_FilterNoMatch(t *testing.T) {
	c, _ := New(makeSource())
	ctx, err := c.Clone("dev", "empty-clone", []string{"NONEXISTENT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ctx.Entries))
	}
}
