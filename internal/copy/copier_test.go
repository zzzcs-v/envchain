package copy

import (
	"testing"
)

func TestCopy_EmptySource(t *testing.T) {
	c := New(false)
	dst := map[string]string{"A": "1"}
	out, res, err := c.Copy(map[string]string{}, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if len(res.Copied) != 0 {
		t.Errorf("expected 0 copied, got %d", len(res.Copied))
	}
}

func TestCopy_NilSource(t *testing.T) {
	c := New(false)
	_, _, err := c.Copy(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestCopy_NoOverwrite_SkipsExisting(t *testing.T) {
	c := New(false)
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}
	out, res, err := c.Copy(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "old" {
		t.Errorf("expected A=old, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	if len(res.Copied) != 1 || res.Copied[0] != "B" {
		t.Errorf("expected B to be copied, got %v", res.Copied)
	}
}

func TestCopy_Overwrite_ReplacesExisting(t *testing.T) {
	c := New(true)
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res, err := c.Copy(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %s", out["A"])
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestCopy_NilDst_InitializesMap(t *testing.T) {
	c := New(false)
	src := map[string]string{"X": "42"}
	out, res, err := c.Copy(src, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "42" {
		t.Errorf("expected X=42, got %s", out["X"])
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 copied")
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{Copied: []string{"A", "B"}, Skipped: []string{"C"}}
	s := r.Summary()
	if s != "copied 2 vars, skipped 1" {
		t.Errorf("unexpected summary: %s", s)
	}
}
