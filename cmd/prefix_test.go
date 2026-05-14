package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envchain/internal/prefix"
)

func runPrefixCmd(t *testing.T, src map[string]string, opts prefix.Options) (string, error) {
	t.Helper()
	out, _, err := prefix.Apply(src, opts)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for k, v := range out {
		buf.WriteString(k + "=" + v + "\n")
	}
	return buf.String(), nil
}

func TestRunPrefix_AddsPrefix(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	out, err := runPrefixCmd(t, src, prefix.Options{Prefix: "TEST_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "TEST_FOO=bar") {
		t.Errorf("expected TEST_FOO=bar in output, got: %s", out)
	}
}

func TestRunPrefix_StripsPrefix(t *testing.T) {
	src := map[string]string{"TEST_FOO": "bar", "OTHER": "val"}
	out, err := runPrefixCmd(t, src, prefix.Options{Prefix: "TEST_", Strip: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "OTHER=val") {
		t.Errorf("expected OTHER=val in output, got: %s", out)
	}
}

func TestRunPrefix_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := runPrefixCmd(t, map[string]string{"K": "v"}, prefix.Options{Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestRunPrefix_CobraCommand_MissingFlag(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	cmd := &cobra.Command{
		Use:  "prefix",
		RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	var pfx string
	cmd.Flags().StringVarP(&pfx, "prefix", "p", "", "prefix string")
	_ = cmd.MarkFlagRequired("prefix")
	root.AddCommand(cmd)
	root.SetArgs([]string{"prefix"})
	if err := root.Execute(); err == nil {
		t.Error("expected error when --prefix flag is missing")
	}
}
