package normalize

import (
	"testing"
)

func TestNormalize_NilSource(t *testing.T) {
	out, err := Normalize(nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestNormalize_NoOptions_ReturnsCopy(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Normalize(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalize_UpperCase(t *testing.T) {
	src := map[string]string{"foo": "1", "bar": "2"}
	out, err := Normalize(src, Options{UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalize_LowerCase(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	out, err := Normalize(src, Options{LowerCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foo"] != "1" || out["bar"] != "2" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalize_ReplacePattern(t *testing.T) {
	src := map[string]string{"foo-bar": "v", "baz.qux": "w"}
	out, err := Normalize(src, Options{ReplacePattern: `[-.]`, ReplaceWith: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foo_bar"] != "v" || out["baz_qux"] != "w" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalize_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := Normalize(map[string]string{"k": "v"}, Options{ReplacePattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNormalize_StripNonAlnum(t *testing.T) {
	src := map[string]string{"foo-bar!": "1"}
	out, err := Normalize(src, Options{StripNonAlnum: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foobar"] != "1" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalize_RestrictToKeys(t *testing.T) {
	src := map[string]string{"foo": "1", "bar": "2"}
	out, err := Normalize(src, Options{UpperCase: true, RestrictToKeys: []string{"foo"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "1" {
		t.Fatalf("expected FOO key, got %v", out)
	}
	if out["bar"] != "2" {
		t.Fatalf("expected bar to be unchanged, got %v", out)
	}
}

func TestNormalize_CollisionReturnsError(t *testing.T) {
	src := map[string]string{"foo": "1", "FOO": "2"}
	_, err := Normalize(src, Options{LowerCase: true})
	if err == nil {
		t.Fatal("expected collision error")
	}
}
