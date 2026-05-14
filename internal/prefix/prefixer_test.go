package prefix

import (
	"testing"
)

func TestApply_EmptyPrefix(t *testing.T) {
	_, _, err := Apply(map[string]string{"KEY": "val"}, Options{Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestApply_NilSource(t *testing.T) {
	out, res, err := Apply(nil, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
	if res.Added != 0 || res.Skipped != 0 {
		t.Errorf("expected zero result, got %+v", res)
	}
}

func TestApply_AddsPrefix(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	out, res, err := Apply(src, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("expected 2 added, got %d", res.Added)
	}
	if out["APP_FOO"] != "1" || out["APP_BAR"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApply_SkipsAlreadyPrefixed(t *testing.T) {
	src := map[string]string{"APP_FOO": "1", "BAR": "2"}
	out, res, err := Apply(src, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 || res.Added != 1 {
		t.Errorf("expected 1 skipped 1 added, got %+v", res)
	}
	if _, ok := out["APP_APP_FOO"]; ok {
		t.Error("should not double-prefix key")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	src := map[string]string{"APP_FOO": "1", "BAR": "2"}
	out, res, err := Apply(src, Options{Prefix: "APP_", Strip: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stripped != 1 || res.Skipped != 1 {
		t.Errorf("expected 1 stripped 1 skipped, got %+v", res)
	}
	if out["FOO"] != "1" {
		t.Errorf("expected FOO=1, got %v", out)
	}
	if out["BAR"] != "2" {
		t.Errorf("expected BAR=2, got %v", out)
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{Added: 3, Skipped: 1}
	if s := r.Summary(); s != "prefixed 3 key(s), skipped 1" {
		t.Errorf("unexpected summary: %s", s)
	}
	r2 := Result{Stripped: 2, Skipped: 0}
	if s := r2.Summary(); s != "stripped 2 key(s), skipped 0" {
		t.Errorf("unexpected summary: %s", s)
	}
}
