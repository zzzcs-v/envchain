package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envchain-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
version: "1"
contexts:
  dev:
    vars:
      APP_ENV: development
  prod:
    extends: dev
    vars:
      APP_ENV: production
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Contexts) != 2 {
		t.Errorf("expected 2 contexts, got %d", len(cfg.Contexts))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	path := writeTemp(t, `contexts:\n  dev:\n    vars:\n      X: y\n`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing version")
	}
}

func TestLoad_DuplicateContext(t *testing.T) {
	// YAML does not produce duplicate keys natively; validate catches unknown extends.
	path := writeTemp(t, `
version: "1"
contexts:
  staging:
    extends: base
    vars:
      APP_ENV: staging
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for extends pointing to unknown context")
	}
}

func TestToResolverInputs(t *testing.T) {
	path := writeTemp(t, `
version: "1"
contexts:
  base:
    vars:
      LOG: info
  dev:
    extends: base
    vars:
      APP_ENV: dev
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defs, extends := cfg.ToResolverInputs()
	if _, ok := defs["base"]; !ok {
		t.Error("expected base in defs")
	}
	if extends["dev"] != "base" {
		t.Errorf("expected dev extends base, got %q", extends["dev"])
	}
}
