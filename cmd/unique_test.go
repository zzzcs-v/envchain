package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

// readCloser wraps a reader to satisfy interfaces used by runUnique.
type readCloser struct{ io.Reader }

func (readCloser) Close() error { return nil }

func runUniqueCmd(input string, keys []string, caseSensitive, showDropped bool) (map[string]string, error) {
	r := readCloser{strings.NewReader(input)}
	// Capture stdout.
	old := bytes.Buffer{}
	_ = old

	// Use direct call with a pipe to capture output.
	var result map[string]string
	pr, pw, _ := func() (io.Reader, io.WriteCloser, error) {
		var buf bytes.Buffer
		return &buf, nopWriteCloser{&buf}, nil
	}()
	_ = pr
	_ = pw

	// Simpler: decode the output by calling the internal package directly.
	import_unique := func() (map[string]string, map[string]string, error) {
		var src map[string]string
		if err := json.NewDecoder(r).Decode(&src); err != nil {
			return nil, nil, err
		}
		import "envchain/internal/unique"
		res, err := unique.Unique(src, unique.Options{Keys: keys, CaseSensitive: caseSensitive})
		return res.Kept, res.Dropped, err
	}
	kept, dropped, err := import_unique()
	if err != nil {
		return nil, err
	}
	if showDropped {
		result = dropped
	} else {
		result = kept
	}
	return result, nil
}

func TestRunUnique_NoDuplicates(t *testing.T) {
	input := `{"A":"foo","B":"bar"}`
	out, err := runUniqueCmd(input, nil, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestRunUnique_DropsDuplicate(t *testing.T) {
	input := `{"A":"same","B":"same","C":"other"}`
	out, err := runUniqueCmd(input, nil, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 kept, got %d", len(out))
	}
}

func TestRunUnique_ShowDropped(t *testing.T) {
	input := `{"A":"dup","B":"dup"}`
	out, err := runUniqueCmd(input, nil, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 dropped, got %d", len(out))
	}
}

func TestRunUnique_InvalidJSON(t *testing.T) {
	_, err := runUniqueCmd(`not-json`, nil, false, false)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
