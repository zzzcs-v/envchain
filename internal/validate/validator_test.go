package validate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/validate"
)

func TestVars_AllValid(t *testing.T) {
	res := validate.Vars(map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	})
	if !res.OK() {
		t.Fatalf("expected no errors, got: %v", res.Errors)
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("expected no warnings, got: %v", res.Warnings)
	}
}

func TestVars_InvalidKey(t *testing.T) {
	res := validate.Vars(map[string]string{
		"invalid-key": "value",
	})
	if res.OK() {
		t.Fatal("expected error for invalid key, got none")
	}
	if !strings.Contains(res.Errors[0], "invalid-key") {
		t.Errorf("error message should mention the bad key, got: %s", res.Errors[0])
	}
}

func TestVars_EmptyValue(t *testing.T) {
	res := validate.Vars(map[string]string{
		"MY_VAR": "",
	})
	if !res.OK() {
		t.Fatalf("empty value should be a warning, not an error")
	}
	if len(res.Warnings) == 0 {
		t.Fatal("expected a warning for empty value")
	}
}

func TestVars_LowercaseKey(t *testing.T) {
	res := validate.Vars(map[string]string{
		"my_var": "hello",
	})
	if res.OK() {
		t.Fatal("lowercase key should produce an error")
	}
}

func TestResult_Summary(t *testing.T) {
	res := validate.Vars(map[string]string{
		"bad-key": "",
	})
	summary := res.Summary()
	if !strings.Contains(summary, "ERROR:") {
		t.Errorf("summary should contain ERROR label, got: %s", summary)
	}
}
