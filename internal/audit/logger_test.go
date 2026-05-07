package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestLog_WritesEntry(t *testing.T) {
	path := tempLog(t)
	l, err := NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer l.Close()

	if err := l.Log("export", "production", []string{"DB_URL", "API_KEY"}, nil); err != nil {
		t.Fatalf("Log: %v", err)
	}
}

func TestReadAll_ParsesEntries(t *testing.T) {
	path := tempLog(t)
	l, err := NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	_ = l.Log("export", "staging", []string{"FOO"}, map[string]string{"format": "dotenv"})
	_ = l.Log("validate", "dev", nil, nil)
	l.Close()

	entries, err := ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Action != "export" {
		t.Errorf("expected action=export, got %s", entries[0].Action)
	}
	if entries[0].Context != "staging" {
		t.Errorf("expected context=staging, got %s", entries[0].Context)
	}
	if entries[0].Meta["format"] != "dotenv" {
		t.Errorf("expected meta format=dotenv")
	}
	if entries[1].Action != "validate" {
		t.Errorf("expected action=validate, got %s", entries[1].Action)
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	_, err := ReadAll("/nonexistent/audit.log")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewLogger_InvalidPath(t *testing.T) {
	_, err := NewLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLog_EmptyContext(t *testing.T) {
	path := tempLog(t)
	l, err := NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer l.Close()

	if err := l.Log("diff", "", nil, nil); err != nil {
		t.Fatalf("Log with empty context: %v", err)
	}

	entries, err := ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 || entries[0].Context != "" {
		t.Errorf("unexpected entries: %+v", entries)
	}
	os.Remove(path)
}
