package sort

import (
	"testing"
)

func TestSort_NilSource(t *testing.T) {
	res, err := Sort(nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pairs) != 0 {
		t.Errorf("expected empty pairs, got %d", len(res.Pairs))
	}
}

func TestSort_ByKey_Asc(t *testing.T) {
	src := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	res, err := Sort(src, Options{By: "key", Order: Asc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, p := range res.Pairs {
		if p.Key != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], p.Key)
		}
	}
}

func TestSort_ByKey_Desc(t *testing.T) {
	src := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	res, err := Sort(src, Options{By: "key", Order: Desc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, p := range res.Pairs {
		if p.Key != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], p.Key)
		}
	}
}

func TestSort_ByValue_Asc(t *testing.T) {
	src := map[string]string{"A": "zebra", "B": "apple", "C": "mango"}
	res, err := Sort(src, Options{By: "value", Order: Asc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"apple", "mango", "zebra"}
	for i, p := range res.Pairs {
		if p.Value != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], p.Value)
		}
	}
}

func TestSort_InvalidField(t *testing.T) {
	_, err := Sort(map[string]string{"A": "1"}, Options{By: "timestamp"})
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}
}

func TestSort_DefaultsToKeyAsc(t *testing.T) {
	src := map[string]string{"Z": "last", "A": "first"}
	res, err := Sort(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pairs[0].Key != "A" {
		t.Errorf("expected first key A, got %q", res.Pairs[0].Key)
	}
}

func TestResult_ToMap(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, _ := Sort(src, Options{})
	m := res.ToMap()
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("ToMap returned unexpected values: %v", m)
	}
}
