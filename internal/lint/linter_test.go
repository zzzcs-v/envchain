package lint

import (
	"strings"
	"testing"
)

func TestRun_NoIssues(t *testing.T) {
	ctxs := map[string]map[string]string{
		"prod": {
			"DATABASE_URL": "postgres://localhost/prod",
			"API_KEY":      "abc123",
		},
	}
	res := Run(ctxs)
	if len(res.Issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(res.Issues), res.Issues)
	}
	if res.HasErrors() {
		t.Error("expected HasErrors to be false")
	}
	if res.Summary() != "no issues found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestRun_LowercaseKey(t *testing.T) {
	ctxs := map[string]map[string]string{
		"dev": {
			"database_url": "postgres://localhost/dev",
		},
	}
	res := Run(ctxs)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if res.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", res.Issues[0].Severity)
	}
	if !strings.Contains(res.Issues[0].Message, "uppercase") {
		t.Errorf("expected message about uppercase, got: %s", res.Issues[0].Message)
	}
}

func TestRun_UnresolvedPlaceholder(t *testing.T) {
	ctxs := map[string]map[string]string{
		"staging": {
			"API_URL": "https://${HOST}/api",
		},
	}
	res := Run(ctxs)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if !strings.Contains(res.Issues[0].Message, "unresolved placeholder") {
		t.Errorf("unexpected message: %s", res.Issues[0].Message)
	}
}

func TestRun_EmptyContextName(t *testing.T) {
	ctxs := map[string]map[string]string{
		"  ": {
			"FOO": "bar",
		},
	}
	res := Run(ctxs)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if res.Issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", res.Issues[0].Severity)
	}
	if !res.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
}

func TestIssue_String(t *testing.T) {
	i := Issue{Context: "prod", Key: "DB", Message: "something wrong", Severity: "error"}
	s := i.String()
	if !strings.Contains(s, "[ERROR]") {
		t.Errorf("expected [ERROR] in string, got: %s", s)
	}
	if !strings.Contains(s, "prod.DB") {
		t.Errorf("expected context.key in string, got: %s", s)
	}
}

func TestResult_Summary_WithIssues(t *testing.T) {
	res := &Result{
		Issues: []Issue{
			{Context: "dev", Key: "X", Message: "bad", Severity: "warning"},
			{Context: "dev", Key: "Y", Message: "worse", Severity: "error"},
		},
	}
	if res.Summary() != "2 issue(s) found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
