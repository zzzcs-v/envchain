package filter_test

import (
	"testing"

	"github.com/envchain/envchain/internal/filter"
)

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	vars := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	got, err := filter.Filter(vars, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestFilter_KeyPrefix(t *testing.T) {
	vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	got, err := filter.Filter(vars, filter.Options{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["DB_URL"]; ok {
		t.Error("DB_URL should have been filtered out")
	}
}

func TestFilter_KeyPattern(t *testing.T) {
	vars := map[string]string{"SECRET_KEY": "abc", "SECRET_TOKEN": "xyz", "PUBLIC_URL": "http://"}
	got, err := filter.Filter(vars, filter.Options{KeyPattern: "^SECRET_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestFilter_InvalidPattern_ReturnsError(t *testing.T) {
	vars := map[string]string{"FOO": "bar"}
	_, err := filter.Filter(vars, filter.Options{KeyPattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	vars := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	got, err := filter.Filter(vars, filter.Options{ExcludeKeys: []string{"BAR", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
	if got["FOO"] != "1" {
		t.Errorf("expected FOO=1, got %q", got["FOO"])
	}
}

func TestFilter_InvertMatch(t *testing.T) {
	vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	got, err := filter.Filter(vars, filter.Options{KeyPrefix: "APP_", InvertMatch: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
	if _, ok := got["DB_URL"]; !ok {
		t.Error("expected DB_URL in inverted result")
	}
}
