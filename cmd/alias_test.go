package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/alias"
)

func setupAliasStore(t *testing.T) (string, *alias.Store) {
	t.Helper()
	dir, err := os.MkdirTemp("", "alias-cmd-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := alias.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return dir, s
}

func TestRunAliasSet_And_List(t *testing.T) {
	dir, _ := setupAliasStore(t)
	if err := runAliasSet(dir, "dev", "development"); err != nil {
		t.Fatalf("runAliasSet: %v", err)
	}
	// redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_ = runAliasList(dir)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()
	if !strings.Contains(out, "dev") || !strings.Contains(out, "development") {
		t.Errorf("expected list output to contain alias, got: %s", out)
	}
}

func TestRunAliasList_Empty(t *testing.T) {
	dir, _ := setupAliasStore(t)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_ = runAliasList(dir)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	if !strings.Contains(buf.String(), "no aliases") {
		t.Error("expected 'no aliases' message")
	}
}

func TestRunAliasDelete_NotFound(t *testing.T) {
	dir, _ := setupAliasStore(t)
	if err := runAliasDelete(dir, "ghost"); err == nil {
		t.Error("expected error deleting non-existent alias")
	}
}

func TestRunAliasDelete_Valid(t *testing.T) {
	dir, s := setupAliasStore(t)
	_ = s.Set("stg", "staging")
	if err := runAliasDelete(dir, "stg"); err != nil {
		t.Fatalf("runAliasDelete: %v", err)
	}
	_, err := s.Get("stg")
	if err == nil {
		t.Error("expected alias to be deleted")
	}
	fmt.Println("delete ok")
}
