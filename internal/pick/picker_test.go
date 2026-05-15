package pick_test

import (
	"testing"

	"github.com/envchain/envchain/internal/pick"
)

func src() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
		"DB_PASS":  "secret",
		"LOG_LEVEL": "info",
	}
}

func TestPick_ExplicitKeys(t *testing.T) {
	out, err := pick.Pick(src(), pick.Options{Keys: []string{"APP_HOST", "DB_PASS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", out["APP_HOST"])
	}
}

func TestPick_Pattern(t *testing.T) {
	out, err := pick.Pick(src(), pick.Options{Pattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestPick_KeysAndPattern(t *testing.T) {
	out, err := pick.Pick(src(), pick.Options{Keys: []string{"LOG_LEVEL"}, Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestPick_NoOptions_ReturnsError(t *testing.T) {
	_, err := pick.Pick(src(), pick.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestPick_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := pick.Pick(src(), pick.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestPick_NilSource(t *testing.T) {
	out, err := pick.Pick(nil, pick.Options{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map for nil source, got %d keys", len(out))
	}
}
