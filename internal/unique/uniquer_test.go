package unique

import (
	"testing"
)

func TestUnique_NilSource(t *testing.T) {
	res, err := Unique(nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Kept) != 0 || len(res.Dropped) != 0 {
		t.Fatalf("expected empty result, got kept=%v dropped=%v", res.Kept, res.Dropped)
	}
}

func TestUnique_NoDuplicates(t *testing.T) {
	src := map[string]string{"A": "foo", "B": "bar", "C": "baz"}
	res, err := Unique(src, Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Kept) != 3 {
		t.Fatalf("expected 3 kept, got %d", len(res.Kept))
	}
	if len(res.Dropped) != 0 {
		t.Fatalf("expected 0 dropped, got %d", len(res.Dropped))
	}
}

func TestUnique_DropsDuplicateValue(t *testing.T) {
	src := map[string]string{"A": "same", "B": "same", "C": "other"}
	res, err := Unique(src, Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Kept) != 2 {
		t.Fatalf("expected 2 kept, got %d: %v", len(res.Kept), res.Kept)
	}
	if len(res.Dropped) != 1 {
		t.Fatalf("expected 1 dropped, got %d: %v", len(res.Dropped), res.Dropped)
	}
}

func TestUnique_CaseInsensitive(t *testing.T) {
	src := map[string]string{"A": "Hello", "B": "hello", "C": "world"}
	res, err := Unique(src, Options{CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Dropped) != 1 {
		t.Fatalf("expected 1 dropped (case-insensitive), got %d", len(res.Dropped))
	}
}

func TestUnique_CaseSensitive_NoDrop(t *testing.T) {
	src := map[string]string{"A": "Hello", "B": "hello"}
	res, err := Unique(src, Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Dropped) != 0 {
		t.Fatalf("expected 0 dropped for case-sensitive, got %d", len(res.Dropped))
	}
}

func TestUnique_RestrictToKeys(t *testing.T) {
	src := map[string]string{"A": "dup", "B": "dup", "C": "dup"}
	// Only check A and B; C should always be kept even though value is same.
	res, err := Unique(src, Options{Keys: []string{"A", "B"}, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Kept["C"]; !ok {
		t.Fatal("expected C to be kept (not in restriction set)")
	}
	if len(res.Dropped) != 1 {
		t.Fatalf("expected 1 dropped among restricted keys, got %d", len(res.Dropped))
	}
}

func TestUnique_DeterministicKeepsAlphaFirst(t *testing.T) {
	// A comes before B alphabetically, so A should be kept.
	src := map[string]string{"B": "same", "A": "same"}
	res, _ := Unique(src, Options{CaseSensitive: true})
	if _, ok := res.Kept["A"]; !ok {
		t.Fatal("expected A (alphabetically first) to be kept")
	}
	if _, ok := res.Dropped["B"]; !ok {
		t.Fatal("expected B to be dropped")
	}
}
