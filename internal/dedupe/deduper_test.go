package dedupe_test

import (
	"testing"

	"github.com/envchain/envchain/internal/dedupe"
)

func TestDedupe_NoSources(t *testing.T) {
	res, err := dedupe.Dedupe(nil, dedupe.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty map, got %v", res.Vars)
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	src := []map[string]string{
		{"A": "1", "B": "2"},
		{"C": "3"},
	}
	res, err := dedupe.Dedupe(src, dedupe.Options{Strategy: dedupe.KeepFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", res.Removed)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" || res.Vars["C"] != "3" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	src := []map[string]string{
		{"HOST": "localhost"},
		{"HOST": "remotehost"},
	}
	res, err := dedupe.Dedupe(src, dedupe.Options{Strategy: dedupe.KeepFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", res.Vars["HOST"])
	}
	if len(res.Removed) != 1 || res.Removed[0] != "HOST" {
		t.Errorf("expected HOST in removed, got %v", res.Removed)
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	src := []map[string]string{
		{"HOST": "localhost"},
		{"HOST": "remotehost"},
	}
	res, err := dedupe.Dedupe(src, dedupe.Options{Strategy: dedupe.KeepLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["HOST"] != "remotehost" {
		t.Errorf("expected remotehost, got %q", res.Vars["HOST"])
	}
}

func TestDedupe_RestrictToKeys(t *testing.T) {
	src := []map[string]string{
		{"A": "first", "B": "first"},
		{"A": "second", "B": "second"},
	}
	res, err := dedupe.Dedupe(src, dedupe.Options{
		Strategy:       dedupe.KeepFirst,
		RestrictToKeys: []string{"A"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A should be deduped (keep first), B should use last value (unrestricted)
	if res.Vars["A"] != "first" {
		t.Errorf("expected A=first, got %q", res.Vars["A"])
	}
	// B appears in both maps; without restriction both are processed so last wins
	if res.Vars["B"] == "" {
		t.Errorf("expected B to be set")
	}
	for _, k := range res.Removed {
		if k == "B" {
			t.Errorf("B should not be in removed list")
		}
	}
}

func TestResult_Summary(t *testing.T) {
	r := dedupe.Result{Vars: map[string]string{}, Removed: []string{}}
	if r.Summary() != "no duplicates found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
	r.Removed = []string{"FOO", "BAR"}
	s := r.Summary()
	if s == "" || s == "no duplicates found" {
		t.Errorf("expected non-empty summary with count, got: %s", s)
	}
}
