package rotate

import (
	"errors"
	"strings"
	"testing"
)

func stubGen(val string) GeneratorFunc {
	return func(key string) (string, error) {
		return val + "_" + key, nil
	}
}

func errGen(key string) (string, error) {
	return "", errors.New("gen failed")
}

func TestRotate_NoGenerator(t *testing.T) {
	_, err := Rotate(map[string]string{"A": "1"}, Options{})
	if err == nil {
		t.Fatal("expected error for nil generator")
	}
}

func TestRotate_AllKeys(t *testing.T) {
	vars := map[string]string{"A": "", "B": ""}
	res, err := Rotate(vars, Options{Generator: stubGen("new"), Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Fatalf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if vars["A"] != "new_A" || vars["B"] != "new_B" {
		t.Errorf("unexpected values: %v", vars)
	}
}

func TestRotate_SkipsNonEmpty_WhenNoOverwrite(t *testing.T) {
	vars := map[string]string{"A": "existing", "B": ""}
	res, err := Rotate(vars, Options{Generator: stubGen("new"), Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	if vars["A"] != "existing" {
		t.Error("A should not have been overwritten")
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	vars := map[string]string{"A": "old", "B": "old", "C": "old"}
	res, err := Rotate(vars, Options{Keys: []string{"A", "C"}, Generator: stubGen("v"), Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Fatalf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if vars["B"] != "old" {
		t.Error("B should be unchanged")
	}
}

func TestRotate_MissingKeySkipped(t *testing.T) {
	vars := map[string]string{"A": "val"}
	res, _ := Rotate(vars, Options{Keys: []string{"A", "MISSING"}, Generator: stubGen("x"), Overwrite: true})
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING skipped, got %v", res.Skipped)
	}
}

func TestRotate_GeneratorError(t *testing.T) {
	vars := map[string]string{"A": ""}
	res, err := Rotate(vars, Options{Generator: errGen, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if len(res.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(res.Errors))
	}
	if !strings.Contains(res.Errors[0].Error(), "gen failed") {
		t.Errorf("unexpected error message: %v", res.Errors[0])
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Rotated: []Entry{{Key: "A"}, {Key: "B"}},
		Skipped: []string{"C"},
		Errors:  []error{errors.New("oops")},
	}
	got := r.Summary()
	if got != "rotated=2 skipped=1 errors=1" {
		t.Errorf("unexpected summary: %s", got)
	}
}
