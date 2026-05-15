package rename

import (
	"testing"
)

func TestRename_NilSource(t *testing.T) {
	out, results, err := Rename(nil, Options{FromPattern: ".*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Fatalf("expected empty output for nil source")
	}
}

func TestRename_EmptyPattern(t *testing.T) {
	_, _, err := Rename(map[string]string{"A": "1"}, Options{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestRename_InvalidPattern(t *testing.T) {
	_, _, err := Rename(map[string]string{"A": "1"}, Options{FromPattern: "["})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestRename_NoMatchPassesThrough(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, results, err := Rename(src, Options{FromPattern: "^NOPE", ToTemplate: "X_$0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatal("keys should be unchanged")
	}
}

func TestRename_PrefixKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, results, err := Rename(src, Options{FromPattern: "^DB_", ToTemplate: "PG_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if out["PG_HOST"] != "localhost" {
		t.Errorf("expected PG_HOST=localhost, got %v", out["PG_HOST"])
	}
	if out["PG_PORT"] != "5432" {
		t.Errorf("expected PG_PORT=5432, got %v", out["PG_PORT"])
	}
}

func TestRename_ConflictErrorOnConflict(t *testing.T) {
	src := map[string]string{"OLD_KEY": "v1", "NEW_KEY": "v2"}
	_, _, err := Rename(src, Options{
		FromPattern:     "^OLD_KEY$",
		ToTemplate:      "NEW_KEY",
		ErrorOnConflict: true,
	})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestRename_ConflictSkip(t *testing.T) {
	src := map[string]string{"OLD_KEY": "v1", "NEW_KEY": "v2"}
	out, results, err := Rename(src, Options{
		FromPattern:   "^OLD_KEY$",
		ToTemplate:    "NEW_KEY",
		SkipConflicts: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Error("expected result to be marked skipped")
	}
	if out["OLD_KEY"] != "v1" {
		t.Error("original key should be preserved when skipped")
	}
}
