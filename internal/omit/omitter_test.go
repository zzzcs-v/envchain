package omit

import (
	"testing"
)

func TestOmit_NilSource(t *testing.T) {
	out, err := Omit(nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestOmit_NoOptions_ReturnsCopy(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	out, err := Omit(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestOmit_ExplicitKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, err := Omit(src, Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out["B"] != "2" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestOmit_PatternRemovesMatches(t *testing.T) {
	src := map[string]string{"SECRET_KEY": "x", "SECRET_TOKEN": "y", "PORT": "8080"}
	out, err := Omit(src, Options{Pattern: "^SECRET_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out["PORT"] != "8080" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestOmit_InvalidPattern_ReturnsError(t *testing.T) {
	src := map[string]string{"A": "1"}
	_, err := Omit(src, Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestOmit_PrefixRemovesMatchingKeys(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "3000", "DB_URL": "postgres://"}
	out, err := Omit(src, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out["DB_URL"] != "postgres://" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestOmit_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	_, err := Omit(src, Options{Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["A"]; !ok {
		t.Fatal("source map was mutated")
	}
}
