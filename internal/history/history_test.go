package history

import (
	"os"
	"testing"
	"time"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-history-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestRecord_AndList(t *testing.T) {
	s := tempStore(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.Record("dev", "dotenv", vars); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Context != "dev" {
		t.Errorf("context: got %q, want %q", e.Context, "dev")
	}
	if e.Format != "dotenv" {
		t.Errorf("format: got %q, want %q", e.Format, "dotenv")
	}
	if e.Vars["FOO"] != "bar" {
		t.Errorf("vars FOO: got %q", e.Vars["FOO"])
	}
}

func TestList_SortedByTimestamp(t *testing.T) {
	s := tempStore(t)
	for _, ctx := range []string{"prod", "dev", "staging"} {
		time.Sleep(2 * time.Millisecond)
		if err := s.Record(ctx, "export", nil); err != nil {
			t.Fatalf("Record %s: %v", ctx, err)
		}
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Context != "prod" || entries[2].Context != "staging" {
		t.Errorf("unexpected order: %v", []string{entries[0].Context, entries[1].Context, entries[2].Context})
	}
}

func TestList_Empty(t *testing.T) {
	s := tempStore(t)
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestClear_RemovesAll(t *testing.T) {
	s := tempStore(t)
	_ = s.Record("dev", "json", nil)
	_ = s.Record("prod", "json", nil)
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	entries, _ := s.List()
	if len(entries) != 0 {
		t.Errorf("expected 0 after clear, got %d", len(entries))
	}
}

func TestNewStore_InvalidDir(t *testing.T) {
	_, err := NewStore("/dev/null/bad/path")
	if err == nil {
		t.Error("expected error for invalid dir")
	}
}
