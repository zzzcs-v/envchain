package trim

import (
	"testing"
)

func TestTrim_NoOptions_ReturnsCopy(t *testing.T) {
	src := map[string]string{"A": "hello", "B": "world"}
	out, results, err := Trim(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 || out["A"] != "hello" || out["B"] != "world" {
		t.Errorf("expected unchanged copy, got %v", out)
	}
	for _, r := range results {
		if r.Changed() {
			t.Errorf("expected no changes, got %v", r)
		}
	}
}

func TestTrim_Whitespace(t *testing.T) {
	src := map[string]string{"KEY": "  value  "}
	out, _, err := Trim(src, Options{Whitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["KEY"]; got != "value" {
		t.Errorf("expected %q, got %q", "value", got)
	}
}

func TestTrim_PrefixAndSuffix(t *testing.T) {
	src := map[string]string{"TOKEN": "Bearer abc123;"}
	out, _, err := Trim(src, Options{Prefix: "Bearer ", Suffix: ";"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["TOKEN"]; got != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", got)
	}
}

func TestTrim_RestrictedToKeys(t *testing.T) {
	src := map[string]string{"A": "  hi  ", "B": "  bye  "}
	out, _, err := Trim(src, Options{Keys: []string{"A"}, Whitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["A"]; got != "hi" {
		t.Errorf("A: expected %q, got %q", "hi", got)
	}
	if got := out["B"]; got != "  bye  " {
		t.Errorf("B: expected unchanged, got %q", got)
	}
}

func TestTrim_NilSource(t *testing.T) {
	out, results, err := Trim(nil, Options{Whitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}

func TestResult_String(t *testing.T) {
	changed := Result{Key: "X", Original: "old", Trimmed: "new"}
	if s := changed.String(); s != `X: "old" -> "new"` {
		t.Errorf("unexpected string: %s", s)
	}
	same := Result{Key: "Y", Original: "val", Trimmed: "val"}
	if s := same.String(); s != "Y: unchanged" {
		t.Errorf("unexpected string: %s", s)
	}
}
