package snapshot

import (
	"os"
	"testing"
	"time"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestSave_AndLoad(t *testing.T) {
	store := tempStore(t)
	vars := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}

	if err := store.Save("staging", vars); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := store.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.Context != "staging" {
		t.Errorf("context: got %q, want %q", entry.Context, "staging")
	}
	if entry.Vars["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q, want %q", entry.Vars["DB_HOST"], "localhost")
	}
	if entry.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if time.Since(entry.CreatedAt) > 5*time.Second {
		t.Error("CreatedAt seems too old")
	}
}

func TestLoad_MissingSnapshot(t *testing.T) {
	store := tempStore(t)
	_, err := store.Load("prod")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestSave_Overwrites(t *testing.T) {
	store := tempStore(t)

	_ = store.Save("dev", map[string]string{"KEY": "old"})
	_ = store.Save("dev", map[string]string{"KEY": "new"})

	entry, err := store.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.Vars["KEY"] != "new" {
		t.Errorf("expected overwritten value %q, got %q", "new", entry.Vars["KEY"])
	}
}

func TestDelete_RemovesFile(t *testing.T) {
	store := tempStore(t)
	_ = store.Save("dev", map[string]string{"X": "1"})

	if err := store.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("dev")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestDelete_NonExistent_NoError(t *testing.T) {
	store := tempStore(t)
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("expected no error deleting non-existent snapshot, got: %v", err)
	}
}

func TestNewStore_InvalidDir(t *testing.T) {
	// point at an existing file so MkdirAll fails
	f, _ := os.CreateTemp("", "envchain-test-*")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	_, err := NewStore(f.Name() + "/subdir")
	if err == nil {
		t.Fatal("expected error for invalid dir, got nil")
	}
}
