package context

import (
	"os"
	"testing"
)

func makeResolver(defs map[string]map[string]string, extends map[string]string) *Resolver {
	return NewResolver(defs, extends)
}

func TestResolve_Simple(t *testing.T) {
	r := makeResolver(
		map[string]map[string]string{
			"dev": {"APP_ENV": "development", "PORT": "3000"},
		},
		map[string]string{},
	)
	ctx, err := r.Resolve("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Vars["APP_ENV"] != "development" {
		t.Errorf("expected development, got %s", ctx.Vars["APP_ENV"])
	}
}

func TestResolve_Extends(t *testing.T) {
	r := makeResolver(
		map[string]map[string]string{
			"base":    {"APP_ENV": "base", "LOG_LEVEL": "info"},
			"staging": {"APP_ENV": "staging"},
		},
		map[string]string{"staging": "base"},
	)
	ctx, err := r.Resolve("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Vars["APP_ENV"] != "staging" {
		t.Errorf("expected staging override, got %s", ctx.Vars["APP_ENV"])
	}
	if ctx.Vars["LOG_LEVEL"] != "info" {
		t.Errorf("expected inherited LOG_LEVEL=info, got %s", ctx.Vars["LOG_LEVEL"])
	}
}

func TestResolve_CircularExtends(t *testing.T) {
	r := makeResolver(
		map[string]map[string]string{
			"a": {"X": "1"},
			"b": {"Y": "2"},
		},
		map[string]string{"a": "b", "b": "a"},
	)
	_, err := r.Resolve("a")
	if err == nil {
		t.Fatal("expected circular extends error, got nil")
	}
}

func TestResolve_MissingContext(t *testing.T) {
	r := makeResolver(map[string]map[string]string{}, map[string]string{})
	_, err := r.Resolve("ghost")
	if err == nil {
		t.Fatal("expected error for missing context")
	}
}

func TestResolve_EnvExpansion(t *testing.T) {
	os.Setenv("TEST_SECRET", "supersecret")
	defer os.Unsetenv("TEST_SECRET")

	r := makeResolver(
		map[string]map[string]string{
			"prod": {"DB_PASS": "${TEST_SECRET}"},
		},
		map[string]string{},
	)
	ctx, err := r.Resolve("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Vars["DB_PASS"] != "supersecret" {
		t.Errorf("expected supersecret, got %s", ctx.Vars["DB_PASS"])
	}
}

func TestToEnvSlice(t *testing.T) {
	c := &Context{Name: "dev", Vars: map[string]string{"port": "8080"}}
	slice := c.ToEnvSlice()
	if len(slice) != 1 || slice[0] != "PORT=8080" {
		t.Errorf("unexpected env slice: %v", slice)
	}
}
