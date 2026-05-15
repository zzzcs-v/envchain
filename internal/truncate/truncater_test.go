package truncate

import (
	"testing"
)

func TestTruncate_NilSource(t *testing.T) {
	out, results, err := Truncate(nil, Options{MaxLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Fatal("expected empty output for nil source")
	}
}

func TestTruncate_NegativeMaxLen(t *testing.T) {
	_, _, err := Truncate(map[string]string{"K": "v"}, Options{MaxLen: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxLen")
	}
}

func TestTruncate_ZeroMaxLen_NoOp(t *testing.T) {
	src := map[string]string{"A": "hello", "B": "world"}
	out, results, err := Truncate(src, Options{MaxLen: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatal("expected no truncations when MaxLen is 0")
	}
	if out["A"] != "hello" || out["B"] != "world" {
		t.Fatal("values should be unchanged")
	}
}

func TestTruncate_ShortValues_Unchanged(t *testing.T) {
	src := map[string]string{"KEY": "hi"}
	out, results, err := Truncate(src, Options{MaxLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatal("expected no truncation results")
	}
	if out["KEY"] != "hi" {
		t.Fatalf("expected 'hi', got %q", out["KEY"])
	}
}

func TestTruncate_LongValue_Truncated(t *testing.T) {
	src := map[string]string{"MSG": "hello world"}
	out, results, err := Truncate(src, Options{MaxLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if out["MSG"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out["MSG"])
	}
	if results[0].Original != "hello world" {
		t.Fatalf("unexpected original: %q", results[0].Original)
	}
}

func TestTruncate_WithSuffix(t *testing.T) {
	src := map[string]string{"DESC": "abcdefghij"}
	out, _, err := Truncate(src, Options{MaxLen: 5, Suffix: "..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DESC"] != "abcde..." {
		t.Fatalf("expected 'abcde...', got %q", out["DESC"])
	}
}

func TestTruncate_RestrictedToKeys(t *testing.T) {
	src := map[string]string{"A": "longvalue", "B": "longvalue"}
	out, results, err := Truncate(src, Options{MaxLen: 4, Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "A" {
		t.Fatal("expected only key A to be truncated")
	}
	if out["B"] != "longvalue" {
		t.Fatal("key B should be unchanged")
	}
}
