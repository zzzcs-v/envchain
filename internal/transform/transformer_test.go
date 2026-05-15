package transform

import (
	"testing"
)

func TestApply_NilSource(t *testing.T) {
	res, err := Apply(nil, Options{Op: OpUppercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty map, got %v", res.Vars)
	}
}

func TestApply_Uppercase(t *testing.T) {
	src := map[string]string{"FOO": "hello", "BAR": "world"}
	res, err := Apply(src, Options{Op: OpUppercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "HELLO" || res.Vars["BAR"] != "WORLD" {
		t.Errorf("unexpected result: %v", res.Vars)
	}
	if res.Changed != 2 {
		t.Errorf("expected 2 changed, got %d", res.Changed)
	}
}

func TestApply_Lowercase(t *testing.T) {
	src := map[string]string{"A": "HELLO"}
	res, err := Apply(src, Options{Op: OpLowercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "hello" {
		t.Errorf("expected 'hello', got %q", res.Vars["A"])
	}
}

func TestApply_TrimSpace(t *testing.T) {
	src := map[string]string{"K": "  value  ", "J": "clean"}
	res, err := Apply(src, Options{Op: OpTrimSpace})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["K"] != "value" {
		t.Errorf("expected 'value', got %q", res.Vars["K"])
	}
	if res.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", res.Changed)
	}
}

func TestApply_Base64RoundTrip(t *testing.T) {
	src := map[string]string{"SECRET": "mysecret"}
	enc, err := Apply(src, Options{Op: OpBase64Enc})
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	dec, err := Apply(enc.Vars, Options{Op: OpBase64Dec})
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if dec.Vars["SECRET"] != "mysecret" {
		t.Errorf("round-trip failed: got %q", dec.Vars["SECRET"])
	}
}

func TestApply_RestrictedKeys(t *testing.T) {
	src := map[string]string{"A": "hello", "B": "world"}
	res, err := Apply(src, Options{Op: OpUppercase, Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", res.Vars["A"])
	}
	if res.Vars["B"] != "world" {
		t.Errorf("B should be unchanged, got %q", res.Vars["B"])
	}
	if res.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", res.Changed)
	}
}

func TestApply_UnknownOp(t *testing.T) {
	src := map[string]string{"A": "val"}
	_, err := Apply(src, Options{Op: Op("noop")})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}
