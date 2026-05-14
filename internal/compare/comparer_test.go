package compare

import (
	"testing"
)

func TestCompare_Identical(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	res := Compare(left, right)
	if len(res.Shared) != 2 {
		t.Errorf("expected 2 shared, got %d", len(res.Shared))
	}
	if len(res.Different) != 0 || len(res.OnlyInLeft) != 0 || len(res.OnlyInRight) != 0 {
		t.Error("expected no differences")
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}
	res := Compare(left, right)
	if _, ok := res.OnlyInLeft["B"]; !ok {
		t.Error("expected B to be only in left")
	}
	if len(res.Shared) != 1 {
		t.Errorf("expected 1 shared, got %d", len(res.Shared))
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}
	res := Compare(left, right)
	if _, ok := res.OnlyInRight["C"]; !ok {
		t.Error("expected C to be only in right")
	}
}

func TestCompare_Modified(t *testing.T) {
	left := map[string]string{"A": "old"}
	right := map[string]string{"A": "new"}
	res := Compare(left, right)
	pair, ok := res.Different["A"]
	if !ok {
		t.Fatal("expected A to be different")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompare_Summary(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old"}
	right := map[string]string{"B": "new", "C": "3"}
	res := Compare(left, right)
	s := res.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCompare_Keys_Sorted(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "old"}
	right := map[string]string{"A": "new", "M": "3"}
	res := Compare(left, right)
	keys := res.Keys()
	if len(keys) < 2 {
		t.Fatal("expected at least 2 differing keys")
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := Compare(map[string]string{}, map[string]string{})
	if len(res.Shared) != 0 || len(res.Different) != 0 {
		t.Error("expected empty result for empty maps")
	}
}
