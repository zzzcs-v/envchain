package template

import (
	"strings"
	"testing"
)

func TestRender_NoPlaceholders(t *testing.T) {
	r := NewRenderer(false)
	res, err := r.Render("hello world", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "hello world" {
		t.Errorf("expected 'hello world', got %q", res.Output)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing vars, got %v", res.Missing)
	}
}

func TestRender_SimplePlaceholder(t *testing.T) {
	r := NewRenderer(false)
	env := map[string]string{"APP_ENV": "production"}
	res, err := r.Render("env={{ APP_ENV }}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "env=production" {
		t.Errorf("unexpected output: %q", res.Output)
	}
}

func TestRender_MultiplePlaceholders(t *testing.T) {
	r := NewRenderer(false)
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	res, err := r.Render("{{HOST}}:{{PORT}}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "localhost:5432" {
		t.Errorf("unexpected output: %q", res.Output)
	}
}

func TestRender_MissingVar_NonStrict(t *testing.T) {
	r := NewRenderer(false)
	res, err := r.Render("url={{ DB_URL }}", map[string]string{})
	if err != nil {
		t.Fatalf("expected no error in non-strict mode, got %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_URL" {
		t.Errorf("expected missing [DB_URL], got %v", res.Missing)
	}
	if !strings.Contains(res.Output, "{{ DB_URL }}") {
		t.Errorf("expected placeholder preserved in output, got %q", res.Output)
	}
}

func TestRender_MissingVar_Strict(t *testing.T) {
	r := NewRenderer(true)
	_, err := r.Render("url={{ DB_URL }}", map[string]string{})
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
	if !strings.Contains(err.Error(), "DB_URL") {
		t.Errorf("expected error to mention DB_URL, got: %v", err)
	}
}

func TestRender_DuplicatePlaceholder_ReportedOnce(t *testing.T) {
	r := NewRenderer(false)
	res, err := r.Render("{{ FOO }} and {{ FOO }}", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing var, got %v", res.Missing)
	}
}
