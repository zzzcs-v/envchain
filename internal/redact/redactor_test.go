package redact

import (
	"regexp"
	"testing"
)

func makeRedactor(t *testing.T) *Redactor {
	t.Helper()
	r, err := New(Options{
		Rules: []Rule{
			{Name: "secret", Pattern: regexp.MustCompile(`secret|password|token`)},
			{Name: "key", Pattern: regexp.MustCompile(`_key$`)},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return r
}

func TestNew_NilPattern(t *testing.T) {
	_, err := New(Options{
		Rules: []Rule{{Name: "bad", Pattern: nil}},
	})
	if err == nil {
		t.Fatal("expected error for nil pattern")
	}
}

func TestNew_DefaultPlaceholder(t *testing.T) {
	r, _ := New(Options{})
	if r.opts.Placeholder != "[REDACTED]" {
		t.Errorf("expected default placeholder, got %q", r.opts.Placeholder)
	}
}

func TestRedactMap_RedactsSensitiveKeys(t *testing.T) {
	r := makeRedactor(t)
	input := map[string]string{
		"PASSWORD": "hunter2",
		"API_KEY":  "abc123",
		"HOST":     "localhost",
	}
	out := r.RedactMap(input)
	if out["PASSWORD"] != "[REDACTED]" {
		t.Errorf("PASSWORD not redacted: %q", out["PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY not redacted: %q", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", out["HOST"])
	}
}

func TestRedactMap_DoesNotMutateInput(t *testing.T) {
	r := makeRedactor(t)
	input := map[string]string{"token": "secret-val"}
	r.RedactMap(input)
	if input["token"] != "secret-val" {
		t.Error("original map was mutated")
	}
}

func TestRedactString_ReplacesMatches(t *testing.T) {
	r, _ := New(Options{
		Rules: []Rule{
			{Name: "num", Pattern: regexp.MustCompile(`\d+`), Replace: "***"},
		},
		Placeholder: "[REDACTED]",
	})
	result := r.RedactString("port=8080 and id=42")
	expected := "port=*** and id=***"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRedactMap_CustomReplacement(t *testing.T) {
	r, _ := New(Options{
		Rules: []Rule{
			{Name: "pw", Pattern: regexp.MustCompile(`password`), Replace: "<hidden>"},
		},
	})
	out := r.RedactMap(map[string]string{"password": "s3cr3t", "user": "admin"})
	if out["password"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", out["password"])
	}
	if out["user"] != "admin" {
		t.Errorf("user should be unchanged")
	}
}
