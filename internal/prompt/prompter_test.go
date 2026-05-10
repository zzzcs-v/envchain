package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func makePrompter(input string) (*Prompter, *bytes.Buffer) {
	out := &bytes.Buffer{}
	p := NewWithIO(strings.NewReader(input), out)
	return p, out
}

func TestAsk_ReturnsInput(t *testing.T) {
	p, out := makePrompter("my-value\n")
	val, err := p.Ask("Enter key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "my-value" {
		t.Errorf("expected %q, got %q", "my-value", val)
	}
	if !strings.Contains(out.String(), "Enter key") {
		t.Error("expected label in output")
	}
}

func TestAsk_TrimsWhitespace(t *testing.T) {
	p, _ := makePrompter("  spaced  \n")
	val, err := p.Ask("Label")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "spaced" {
		t.Errorf("expected trimmed value, got %q", val)
	}
}

func TestAsk_NoInput_ReturnsError(t *testing.T) {
	p, _ := makePrompter("")
	_, err := p.Ask("Label")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestAskSecret_FallbackReturnsInput(t *testing.T) {
	// Non-terminal fallback path is exercised here.
	p, _ := makePrompter("s3cr3t\n")
	val, err := p.AskSecret("Passphrase")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected %q, got %q", "s3cr3t", val)
	}
}

func TestConfirm_Yes(t *testing.T) {
	for _, input := range []string{"y\n", "Y\n", "yes\n", "YES\n"} {
		p, _ := makePrompter(input)
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Errorf("expected true for input %q", input)
		}
	}
}

func TestConfirm_No(t *testing.T) {
	for _, input := range []string{"n\n", "no\n", "\n", "maybe\n"} {
		p, _ := makePrompter(input)
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Errorf("expected false for input %q", input)
		}
	}
}

func TestConfirm_OutputContainsLabel(t *testing.T) {
	p, out := makePrompter("y\n")
	_, _ = p.Confirm("Delete everything?")
	if !strings.Contains(out.String(), "Delete everything?") {
		t.Error("expected label in confirm output")
	}
}
