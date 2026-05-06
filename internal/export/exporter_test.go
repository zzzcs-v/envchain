package export

import (
	"bytes"
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_URL":   "postgres://localhost/mydb",
	"SECRET":   "hunter2",
	"WITH_SPACE": "hello world",
}

func TestNewExporter_InvalidFormat(t *testing.T) {
	_, err := NewExporter("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_Dotenv(t *testing.T) {
	ex, _ := NewExporter(FormatDotenv)
	var buf bytes.Buffer
	if err := ex.Write(&buf, sampleEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("missing APP_ENV line in dotenv output")
	}
	if !strings.Contains(out, `WITH_SPACE="hello world"`) {
		t.Errorf("expected quoted value for WITH_SPACE, got:\n%s", out)
	}
}

func TestWrite_Export(t *testing.T) {
	ex, _ := NewExporter(FormatExport)
	var buf bytes.Buffer
	if err := ex.Write(&buf, sampleEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export DB_URL=postgres://localhost/mydb") {
		t.Errorf("missing export prefix in output:\n%s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	ex, _ := NewExporter(FormatJSON)
	var buf bytes.Buffer
	if err := ex.Write(&buf, sampleEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("output is not valid JSON-like structure:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("missing APP_ENV in JSON output:\n%s", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	ex, _ := NewExporter(FormatDotenv)
	var buf bytes.Buffer
	_ = ex.Write(&buf, sampleEnv)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output is not sorted: %q comes after %q", lines[i], lines[i-1])
		}
	}
}
