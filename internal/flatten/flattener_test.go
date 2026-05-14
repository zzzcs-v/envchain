package flatten

import (
	"testing"
)

func TestFlatten_NilSource(t *testing.T) {
	res, err := Flatten(nil, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty map, got %v", res.Vars)
	}
}

func TestFlatten_FlatMap(t *testing.T) {
	src := map[string]any{"FOO": "bar", "BAZ": "qux"}
	res, err := Flatten(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "bar" || res.Vars["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	src := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	res, err := Flatten(src, Options{Separator: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %v", res.Vars)
	}
	if res.Vars["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %v", res.Vars)
	}
}

func TestFlatten_UpperCase(t *testing.T) {
	src := map[string]any{"app": map[string]any{"name": "envchain"}}
	res, err := Flatten(src, Options{UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_NAME"] != "envchain" {
		t.Errorf("expected APP_NAME=envchain, got %v", res.Vars)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	src := map[string]any{"KEY": "val"}
	res, err := Flatten(src, Options{Prefix: "MY", Separator: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["MY_KEY"] != "val" {
		t.Errorf("expected MY_KEY=val, got %v", res.Vars)
	}
}

func TestFlatten_UnsupportedType_EmitsWarning(t *testing.T) {
	src := map[string]any{"LIST": []string{"a", "b"}}
	res, err := Flatten(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for slice type, got none")
	}
	if _, ok := res.Vars["LIST"]; ok {
		t.Error("slice key should have been skipped")
	}
}

func TestFlatten_NilValue_EmptyString(t *testing.T) {
	src := map[string]any{"EMPTY": nil}
	res, err := Flatten(src, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := res.Vars["EMPTY"]; !ok || v != "" {
		t.Errorf("expected EMPTY=\"\", got %q", v)
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	src := map[string]any{"a": map[string]any{"b": "c"}}
	res, err := Flatten(src, Options{Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["a.b"] != "c" {
		t.Errorf("expected a.b=c, got %v", res.Vars)
	}
}
