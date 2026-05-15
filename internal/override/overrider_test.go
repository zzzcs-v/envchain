package override

import (
	"testing"
)

func TestParsePair_Valid(t *testing.T) {
	e, err := ParsePair("FOO=bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Key != "FOO" || e.Value != "bar" {
		t.Fatalf("got %+v", e)
	}
}

func TestParsePair_EmptyValue(t *testing.T) {
	e, err := ParsePair("FOO=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Key != "FOO" || e.Value != "" {
		t.Fatalf("got %+v", e)
	}
}

func TestParsePair_Invalid(t *testing.T) {
	for _, s := range []string{"", "NOEQUALS", "=VAL"} {
		_, err := ParsePair(s)
		if err == nil {
			t.Fatalf("expected error for %q", s)
		}
	}
}

func TestApply_BasicOverride(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	err := Apply(dst, []Entry{{Key: "A", Value: "99"}}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if dst["A"] != "99" {
		t.Fatalf("expected 99, got %s", dst["A"])
	}
}

func TestApply_NilDst(t *testing.T) {
	err := Apply(nil, []Entry{{Key: "A", Value: "1"}}, Options{})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestApply_EmptyValueBlocked(t *testing.T) {
	dst := map[string]string{"A": "1"}
	err := Apply(dst, []Entry{{Key: "A", Value: ""}}, Options{AllowEmpty: false})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestApply_EmptyValueAllowed(t *testing.T) {
	dst := map[string]string{"A": "1"}
	err := Apply(dst, []Entry{{Key: "A", Value: ""}}, Options{AllowEmpty: true})
	if err != nil {
		t.Fatal(err)
	}
	if dst["A"] != "" {
		t.Fatal("expected empty string")
	}
}

func TestApply_StrictMode_MissingKey(t *testing.T) {
	dst := map[string]string{"A": "1"}
	err := Apply(dst, []Entry{{Key: "Z", Value: "9"}}, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestApply_StrictMode_ExistingKey(t *testing.T) {
	dst := map[string]string{"A": "1"}
	err := Apply(dst, []Entry{{Key: "A", Value: "2"}}, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
}
