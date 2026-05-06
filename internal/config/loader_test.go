package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envchain.json")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `{
		"version": "1",
		"contexts": [
			{"name": "Development", "context": "dev", "vars": {"API_URL": "http://localhost"}},
			{"name": "Production", "context": "prod", "vars": {"API_URL": "https://api.example.com"}}
		]
	}`
	path := writeTemp(t, raw)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Contexts) != 2 {
		t.Errorf("expected 2 contexts, got %d", len(cfg.Contexts))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/.envchain.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	raw := `{"contexts": [{"context": "dev", "vars": {}}]}`
	path := writeTemp(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing version")
	}
}

func TestLoad_DuplicateContext(t *testing.T) {
	raw := `{
		"version": "1",
		"contexts": [
			{"context": "dev", "vars": {}},
			{"context": "dev", "vars": {}}
		]
	}`
	path := writeTemp(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate context names")
	}
}

func TestGetContext(t *testing.T) {
	cfg := &ChainConfig{
		Version: "1",
		Contexts: []EnvConfig{
			{Context: "dev", Vars: map[string]string{"FOO": "bar"}},
		},
	}
	ctx, err := cfg.GetContext("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", ctx.Vars["FOO"])
	}
	_, err = cfg.GetContext("staging")
	if err == nil {
		t.Fatal("expected error for missing context")
	}
}
