package squash

import (
	"testing"
)

func TestSquash_NoSources(t *testing.T) {
	res, err := Squash(nil, Options{Strategy: KeepFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty map, got %v", res.Vars)
	}
}

func TestSquash_SingleSource(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	res, err := Squash([]map[string]string{src}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestSquash_KeepFirst(t *testing.T) {
	a := map[string]string{"X": "original"}
	b := map[string]string{"X": "override"}
	res, err := Squash([]map[string]string{a, b}, Options{Strategy: KeepFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["X"] != "original" {
		t.Errorf("expected 'original', got %q", res.Vars["X"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "X" {
		t.Errorf("expected conflict on X, got %v", res.Conflicts)
	}
}

func TestSquash_KeepLast(t *testing.T) {
	a := map[string]string{"X": "original"}
	b := map[string]string{"X": "override"}
	res, err := Squash([]map[string]string{a, b}, Options{Strategy: KeepLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["X"] != "override" {
		t.Errorf("expected 'override', got %q", res.Vars["X"])
	}
}

func TestSquash_ErrorOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "a"}
	b := map[string]string{"KEY": "b"}
	_, err := Squash([]map[string]string{a, b}, Options{Strategy: ErrorOnConflict})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestSquash_NilSourceSkipped(t *testing.T) {
	a := map[string]string{"A": "1"}
	res, err := Squash([]map[string]string{a, nil}, Options{Strategy: KeepFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "1" {
		t.Errorf("expected A=1, got %v", res.Vars)
	}
}

func TestParseStrategy_Valid(t *testing.T) {
	cases := map[string]Strategy{
		"keep-first": KeepFirst,
		"keep-last":  KeepLast,
		"error":      ErrorOnConflict,
	}
	for input, want := range cases {
		got, err := ParseStrategy(input)
		if err != nil {
			t.Errorf("ParseStrategy(%q) error: %v", input, err)
		}
		if got != want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := ParseStrategy("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
