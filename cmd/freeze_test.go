package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func runFreezeCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"freeze"}, args...))
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

func TestRunFreeze_SnapshotOutput(t *testing.T) {
	out, err := runFreezeCmd("--set", "APP_ENV=prod", "--set", "PORT=8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV") || !strings.Contains(out, "prod") {
		t.Fatalf("expected snapshot JSON in output, got: %s", out)
	}
}

func TestRunFreeze_NoDrift(t *testing.T) {
	out, err := runFreezeCmd(
		"--set", "APP_ENV=prod",
		"--diff", "APP_ENV=prod",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no drift") {
		t.Fatalf("expected no drift message, got: %s", out)
	}
}

func TestRunFreeze_DetectsDrift(t *testing.T) {
	out, err := runFreezeCmd(
		"--set", "PORT=8080",
		"--diff", "PORT=9000",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "PORT") {
		t.Fatalf("expected PORT in drift output, got: %s", out)
	}
}

func TestRunFreeze_InvalidSetPair(t *testing.T) {
	_, err := runFreezeCmd("--set", "BADPAIR")
	if err == nil {
		t.Fatal("expected error for invalid --set pair")
	}
}

func TestRunFreeze_InvalidDiffPair(t *testing.T) {
	_, err := runFreezeCmd("--set", "KEY=val", "--diff", "NODIFF")
	if err == nil {
		t.Fatal("expected error for invalid --diff pair")
	}
}
