package schema

import (
	"strings"
	"testing"
)

func makeSchema(t *testing.T, rules map[string]Rule) *Schema {
	t.Helper()
	s, err := NewSchema(rules)
	if err != nil {
		t.Fatalf("NewSchema failed: %v", err)
	}
	return s
}

func TestValidate_NoViolations(t *testing.T) {
	s := makeSchema(t, map[string]Rule{
		"PORT": {Required: true, Pattern: `^\d+$`},
	})
	violations := s.Validate(map[string]string{"PORT": "8080"})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	s := makeSchema(t, map[string]Rule{
		"DB_URL": {Required: true},
	})
	violations := s.Validate(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "required") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := makeSchema(t, map[string]Rule{
		"PORT": {Pattern: `^\d+$`},
	})
	violations := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	s := makeSchema(t, map[string]Rule{
		"API_KEY": {Required: true},
	})
	violations := s.Validate(map[string]string{"API_KEY": "   "})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "empty") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestNewSchema_InvalidPattern(t *testing.T) {
	_, err := NewSchema(map[string]Rule{
		"BAD": {Pattern: `[invalid`},
	})
	if err == nil {
		t.Error("expected error for invalid pattern, got nil")
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Key: "FOO", Message: "some issue"}
	if v.String() != "FOO: some issue" {
		t.Errorf("unexpected string: %s", v.String())
	}
}
