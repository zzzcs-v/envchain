package freeze_test

import (
	"testing"

	"github.com/envchain/envchain/internal/freeze"
)

func makeSource() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
}

func TestNew_NilSource(t *testing.T) {
	_, err := freeze.New(nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestNew_ValidSource(t *testing.T) {
	f, err := freeze.New(makeSource())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Freezer")
	}
}

func TestGet_ExistingKey(t *testing.T) {
	f, _ := freeze.New(makeSource())
	v, ok := f.Get("PORT")
	if !ok || v != "8080" {
		t.Fatalf("expected PORT=8080, got %q ok=%v", v, ok)
	}
}

func TestGet_MissingKey(t *testing.T) {
	f, _ := freeze.New(makeSource())
	_, ok := f.Get("MISSING")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestKeys_SortedOutput(t *testing.T) {
	f, _ := freeze.New(makeSource())
	keys := f.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "APP_ENV" || keys[1] != "LOG_LEVEL" || keys[2] != "PORT" {
		t.Fatalf("unexpected key order: %v", keys)
	}
}

func TestSnapshot_IsACopy(t *testing.T) {
	src := makeSource()
	f, _ := freeze.New(src)
	snap := f.Snapshot()
	snap["PORT"] = "9999"
	v, _ := f.Get("PORT")
	if v != "8080" {
		t.Fatal("snapshot mutation should not affect freezer")
	}
}

func TestDiffFrom_NoChanges(t *testing.T) {
	f, _ := freeze.New(makeSource())
	diff := f.DiffFrom(makeSource())
	if len(diff) != 0 {
		t.Fatalf("expected no diff, got %v", diff)
	}
}

func TestDiffFrom_ModifiedKey(t *testing.T) {
	f, _ := freeze.New(makeSource())
	live := makeSource()
	live["PORT"] = "9000"
	diff := f.DiffFrom(live)
	if len(diff) != 1 || diff[0] != "PORT" {
		t.Fatalf("expected [PORT], got %v", diff)
	}
}

func TestDiffFrom_AddedKey(t *testing.T) {
	f, _ := freeze.New(makeSource())
	live := makeSource()
	live["NEW_KEY"] = "value"
	diff := f.DiffFrom(live)
	if len(diff) != 1 || diff[0] != "+NEW_KEY" {
		t.Fatalf("expected [+NEW_KEY], got %v", diff)
	}
}
