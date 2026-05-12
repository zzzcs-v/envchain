package rollback_test

import (
	"os"
	"testing"
	"time"

	"github.com/envchain/envchain/internal/rollback"
)

func tempStore(t *testing.T) *rollback.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := rollback.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSave_AndList(t *testing.T) {
	s := tempStore(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.Save("dev", vars); err != nil {
		t.Fatalf("Save: %v", err)
	}
	entries, err := s.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Context != "dev" {
		t.Errorf("expected context 'dev', got %q", entries[0].Context)
	}
	if entries[0].Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", entries[0].Vars["FOO"])
	}
}

func TestSave_EmptyContext(t *testing.T) {
	s := tempStore(t)
	if err := s.Save("", map[string]string{}); err == nil {
		t.Error("expected error for empty context, got nil")
	}
}

func TestList_SortedNewestFirst(t *testing.T) {
	s := tempStore(t)
	for i := 0; i < 3; i++ {
		if err := s.Save("staging", map[string]string{"I": string(rune('0'+i))}); err != nil {
			t.Fatalf("Save: %v", err)
		}
		time.Sleep(2 * time.Millisecond)
	}
	entries, err := s.List("staging")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.After(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted newest first at index %d", i)
		}
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	s := tempStore(t)
	s.Save("prod", map[string]string{"V": "1"})
	time.Sleep(2 * time.Millisecond)
	s.Save("prod", map[string]string{"V": "2"})
	e, err := s.Latest("prod")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e.Vars["V"] != "2" {
		t.Errorf("expected V=2, got %q", e.Vars["V"])
	}
}

func TestLatest_MissingContext(t *testing.T) {
	s := tempStore(t)
	_, err := s.Latest("ghost")
	if err == nil {
		t.Error("expected error for missing context, got nil")
	}
}

func TestNewStore_InvalidDir(t *testing.T) {
	// point to a file, not a directory
	f, _ := os.CreateTemp("", "envchain-rb-*")
	f.Close()
	defer os.Remove(f.Name())
	_, err := rollback.NewStore(f.Name())
	if err == nil {
		t.Error("expected error when dir is a file, got nil")
	}
}
