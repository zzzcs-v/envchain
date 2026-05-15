package defaults_test

import (
	"testing"

	"github.com/user/envchain/internal/defaults"
)

func TestApply_NilDst(t *testing.T) {
	_, err := defaults.Apply(nil, []defaults.Entry{{Key: "FOO", Value: "bar"}})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	dst := map[string]string{"A": "1"}
	res, err := defaults.Apply(dst, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied)+len(res.Skipped)+len(res.Overrode) != 0 {
		t.Fatal("expected empty result")
	}
}

func TestApply_MissingKey(t *testing.T) {
	dst := map[string]string{}
	res, err := defaults.Apply(dst, []defaults.Entry{{Key: "HOST", Value: "localhost"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", dst["HOST"])
	}
	if len(res.Applied) != 1 || res.Applied[0] != "HOST" {
		t.Errorf("expected HOST in Applied, got %v", res.Applied)
	}
}

func TestApply_ExistingNonEmpty_Skipped(t *testing.T) {
	dst := map[string]string{"PORT": "8080"}
	res, err := defaults.Apply(dst, []defaults.Entry{{Key: "PORT", Value: "3000"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["PORT"] != "8080" {
		t.Errorf("expected PORT unchanged, got %q", dst["PORT"])
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", res.Skipped)
	}
}

func TestApply_EmptyValue_AppliesDefault(t *testing.T) {
	dst := map[string]string{"DB": ""}
	_, err := defaults.Apply(dst, []defaults.Entry{{Key: "DB", Value: "postgres"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["DB"] != "postgres" {
		t.Errorf("expected DB=postgres, got %q", dst["DB"])
	}
}

func TestApply_Override_ReplacesExisting(t *testing.T) {
	dst := map[string]string{"ENV": "dev"}
	res, err := defaults.Apply(dst, []defaults.Entry{{Key: "ENV", Value: "prod", Override: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["ENV"] != "prod" {
		t.Errorf("expected ENV=prod, got %q", dst["ENV"])
	}
	if len(res.Overrode) != 1 {
		t.Errorf("expected 1 overrode, got %v", res.Overrode)
	}
}

func TestApply_SkipsEmptyKey(t *testing.T) {
	dst := map[string]string{}
	res, err := defaults.Apply(dst, []defaults.Entry{{Key: "", Value: "ignored"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst) != 0 {
		t.Errorf("expected empty dst, got %v", dst)
	}
	if len(res.Applied) != 0 {
		t.Errorf("expected nothing applied, got %v", res.Applied)
	}
}

func TestResult_Summary(t *testing.T) {
	r := defaults.Result{
		Applied:  []string{"A", "B"},
		Skipped:  []string{"C"},
		Overrode: []string{},
	}
	got := r.Summary()
	want := "applied=2 skipped=1 overrode=0"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
