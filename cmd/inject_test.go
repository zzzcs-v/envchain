package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runInjectCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "inject", RunE: func(cmd *cobra.Command, _ []string) error {
		var (
			overwrite bool
			prefix    string
			sources   []string
		)
		cmd.Flags().BoolVar(&overwrite, "overwrite", false, "")
		cmd.Flags().StringVar(&prefix, "prefix", "", "")
		cmd.Flags().StringArrayVar(&sources, "source", nil, "")
		_ = cmd.ParseFlags(args)
		return runInject(sources, prefix, overwrite, cmd)
	}}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_ = cmd.Execute()
	return buf.String(), nil
}

func TestRunInject_BasicOutput(t *testing.T) {
	srcs := []string{"base:FOO=bar,BAZ=qux"}
	_, err := runInject(srcs, "", false, newTestCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunInject_InvalidSourceFormat(t *testing.T) {
	_, err := parseSources([]string{"badformat"})
	if err == nil {
		t.Fatal("expected error for bad source format")
	}
}

func TestRunInject_InvalidPair(t *testing.T) {
	_, err := parseSources([]string{"name:NOEQUALSSIGN"})
	if err == nil {
		t.Fatal("expected error for missing = in pair")
	}
}

func TestRunInject_PrefixApplied(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := newTestCmd()
	cmd.SetOut(buf)
	err := runInject([]string{"svc:PORT=9090"}, "SVC_", false, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SVC_PORT") {
		t.Errorf("expected SVC_PORT in output, got: %s", buf.String())
	}
}

func TestRunInject_MultipleSources(t *testing.T) {
	srcs := []string{"a:X=1", "b:Y=2"}
	buf := &bytes.Buffer{}
	cmd := newTestCmd()
	cmd.SetOut(buf)
	err := runInject(srcs, "", false, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "source=a") || !strings.Contains(out, "source=b") {
		t.Errorf("expected both sources in output, got: %s", out)
	}
}

func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	return cmd
}
