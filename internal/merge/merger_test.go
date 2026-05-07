package merge

import (
	"testing"
)

func TestMerge_NoLayers(t *testing.T) {
	m := NewMerger(StrategyOverwrite)
	result, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestMerge_SingleLayer(t *testing.T) {
	m := NewMerger(StrategyOverwrite)
	layer := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result, err := m.Merge(layer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMerge_Overwrite(t *testing.T) {
	m := NewMerger(StrategyOverwrite)
	base := map[string]string{"FOO": "base", "SHARED": "base"}
	override := map[string]string{"SHARED": "override", "BAR": "new"}
	result, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SHARED"] != "override" {
		t.Errorf("expected SHARED=override, got %q", result["SHARED"])
	}
	if result["FOO"] != "base" {
		t.Errorf("expected FOO=base, got %q", result["FOO"])
	}
	if result["BAR"] != "new" {
		t.Errorf("expected BAR=new, got %q", result["BAR"])
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	m := NewMerger(StrategyKeepExisting)
	base := map[string]string{"FOO": "original"}
	override := map[string]string{"FOO": "ignored"}
	result, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %q", result["FOO"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	m := NewMerger(StrategyError)
	base := map[string]string{"FOO": "a"}
	conflict := map[string]string{"FOO": "b"}
	_, err := m.Merge(base, conflict)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestMerge_ErrorOnConflict_NoConflict(t *testing.T) {
	m := NewMerger(StrategyError)
	a := map[string]string{"FOO": "a"}
	b := map[string]string{"BAR": "b"}
	result, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}
