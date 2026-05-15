package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"envchain/internal/override"
)

func TestRunOverride_NoPairs(t *testing.T) {
	err := runOverride(nil, false, false)
	if err == nil {
		t.Fatal("expected error for empty pairs")
	}
}

func TestRunOverride_InvalidPair(t *testing.T) {
	err := runOverride([]string{"NOEQUALS"}, false, false)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParsePair_RoundTrip(t *testing.T) {
	e, err := override.ParsePair("MY_KEY=hello world")
	if err != nil {
		t.Fatal(err)
	}
	if e.Key != "MY_KEY" || e.Value != "hello world" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestApply_OutputContainsOverride(t *testing.T) {
	dst := map[string]string{"FOO": "old", "BAR": "keep"}
	err := override.Apply(dst, []override.Entry{{Key: "FOO", Value: "new"}}, override.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if dst["FOO"] != "new" {
		t.Fatalf("expected new, got %s", dst["FOO"])
	}
	if dst["BAR"] != "keep" {
		t.Fatal("BAR should be unchanged")
	}
}

func TestApply_JSONSerializable(t *testing.T) {
	dst := map[string]string{"A": "1"}
	_ = override.Apply(dst, []override.Entry{{Key: "A", Value: "2"}}, override.Options{})

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(dst); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"A"`) {
		t.Fatal("expected key A in JSON output")
	}
}
