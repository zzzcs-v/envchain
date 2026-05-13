package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runMaskCmd(t *testing.T, extraArgs ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	root := &cobra.Command{Use: "envchain"}
	maskCmd := &cobra.Command{
		Use:  "mask",
		RunE: runMask,
	}
	maskCmd.Flags().StringVar(&maskStrategy, "strategy", "redact", "")
	maskCmd.Flags().IntVar(&maskVisible, "visible", 3, "")
	maskCmd.Flags().StringVar(&maskPattern, "pattern", `(?i)(token|secret)`, "")
	maskCmd.SetOut(&buf)
	root.AddCommand(maskCmd)
	root.SetOut(&buf)
	args := append([]string{"mask"}, extraArgs...)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestRunMask_OutputIsJSON(t *testing.T) {
	out, err := runMaskCmd(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON output, got: %q", out)
	}
}

func TestRunMask_InvalidPattern(t *testing.T) {
	_, err := runMaskCmd(t, "--pattern", "[")
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestRunMask_StrategyFlag(t *testing.T) {
	out, err := runMaskCmd(t, "--strategy", "hash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// output should still be valid JSON
	if !strings.Contains(out, "{") {
		t.Errorf("expected JSON, got %q", out)
	}
}
