package mask

import (
	"strings"
	"testing"
)

func makeMasker(t *testing.T, opts Options) *Masker {
	t.Helper()
	m, err := New(opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestMaskMap_RedactsMatchingKeys(t *testing.T) {
	m := makeMasker(t, DefaultOptions())
	input := map[string]string{
		"API_TOKEN": "super-secret",
		"APP_NAME":  "envchain",
	}
	out := m.MaskMap(input)
	if out["API_TOKEN"] != "***" {
		t.Errorf("expected *** got %q", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "envchain" {
		t.Errorf("expected envchain got %q", out["APP_NAME"])
	}
}

func TestMaskMap_PartialStrategy(t *testing.T) {
	opts := DefaultOptions()
	opts.Strategy = StrategyPartial
	opts.VisibleChars = 2
	m := makeMasker(t, opts)
	out := m.MaskMap(map[string]string{"DB_PASSWORD": "abcdefgh"})
	v := out["DB_PASSWORD"]
	if !strings.HasPrefix(v, "ab") || !strings.HasSuffix(v, "gh") {
		t.Errorf("unexpected partial mask: %q", v)
	}
	if !strings.Contains(v, "****") {
		t.Errorf("expected stars in middle: %q", v)
	}
}

func TestMaskMap_HashStrategy(t *testing.T) {
	opts := DefaultOptions()
	opts.Strategy = StrategyHash
	m := makeMasker(t, opts)
	out := m.MaskMap(map[string]string{"SECRET_KEY": "myvalue"})
	v := out["SECRET_KEY"]
	if !strings.HasPrefix(v, "#") {
		t.Errorf("expected hash prefix, got %q", v)
	}
}

func TestMaskMap_NilKeyPattern_MasksAll(t *testing.T) {
	opts := Options{Strategy: StrategyRedact, KeyPattern: ""}
	m := makeMasker(t, opts)
	out := m.MaskMap(map[string]string{"ANYTHING": "value"})
	if out["ANYTHING"] != "***" {
		t.Errorf("expected all keys masked, got %q", out["ANYTHING"])
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	opts := DefaultOptions()
	opts.KeyPattern = "["
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestPartialMask_ShortValue(t *testing.T) {
	v := partialMask("ab", 3)
	if v != "**" {
		t.Errorf("expected all stars for short value, got %q", v)
	}
}

func TestHashMask_Empty(t *testing.T) {
	v := hashMask("")
	if v != "***" {
		t.Errorf("expected *** for empty, got %q", v)
	}
}
