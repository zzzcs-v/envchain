package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/sort"
)

func runSortCmd(t *testing.T, input string, flags map[string]string) (string, error) {
	t.Helper()
	cmd := &cobra.Command{RunE: runSort}
	cmd.Flags().StringVar(&sortBy, "by", "key", "")
	cmd.Flags().StringVar(&sortOrder, "order", "asc", "")
	for k, v := range flags {
		_ = cmd.Flags().Set(k, v)
	}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetIn(strings.NewReader(input))
	err := cmd.Execute()
	return strings.TrimSpace(out.String()), err
}

func TestRunSort_ByKeyAsc(t *testing.T) {
	input := `{"ZEBRA":"1","APPLE":"2","MANGO":"3"}`
	output, err := runSortCmd(t, input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var pairs []sort.Pair
	if err := json.Unmarshal([]byte(output), &pairs); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if pairs[0].Key != "APPLE" {
		t.Errorf("expected first key APPLE, got %q", pairs[0].Key)
	}
}

func TestRunSort_ByValueDesc(t *testing.T) {
	input := `{"A":"zebra","B":"apple","C":"mango"}`
	output, err := runSortCmd(t, input, map[string]string{"by": "value", "order": "desc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var pairs []sort.Pair
	if err := json.Unmarshal([]byte(output), &pairs); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if pairs[0].Value != "zebra" {
		t.Errorf("expected first value zebra, got %q", pairs[0].Value)
	}
}

func TestRunSort_InvalidField(t *testing.T) {
	input := `{"A":"1"}`
	_, err := runSortCmd(t, input, map[string]string{"by": "timestamp"})
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}
}

func TestRunSort_InvalidJSON(t *testing.T) {
	_, err := runSortCmd(t, `not-json`, nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON input")
	}
}
