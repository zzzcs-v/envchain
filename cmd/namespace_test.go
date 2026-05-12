package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"envchain/internal/namespace"
)

func setupNSStore(t *testing.T) (string, *namespace.Store) {
	t.Helper()
	dir, err := os.MkdirTemp("", "ns-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := namespace.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return dir, s
}

func TestRunNamespaceSave_And_List(t *testing.T) {
	dir, s := setupNSStore(t)
	_ = dir

	err := s.Save(namespace.Entry{Name: "svc", Prefix: "SVC_", Contexts: []string{"dev"}})
	if err != nil {
		t.Fatal(err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "svc" {
		t.Errorf("unexpected list: %v", list)
	}
}

func TestRunNamespaceDelete_NotFound(t *testing.T) {
	_, s := setupNSStore(t)
	if err := s.Delete("nonexistent"); err == nil {
		t.Error("expected error deleting missing namespace")
	}
}

func TestRunNamespaceDelete_Valid(t *testing.T) {
	_, s := setupNSStore(t)
	_ = s.Save(namespace.Entry{Name: "gone", Prefix: "G_", Contexts: []string{"prod"}})
	if err := s.Delete("gone"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunNamespaceList_Output(t *testing.T) {
	_, s := setupNSStore(t)
	_ = s.Save(namespace.Entry{Name: "infra", Prefix: "INFRA_", Contexts: []string{"dev", "prod"}})

	var buf bytes.Buffer
	list, _ := s.List()
	for _, e := range list {
		buf.WriteString(e.Name + " " + e.Prefix + " " + strings.Join(e.Contexts, ",") + "\n")
	}
	out := buf.String()
	if !strings.Contains(out, "infra") {
		t.Errorf("expected 'infra' in output, got: %s", out)
	}
	if !strings.Contains(out, "INFRA_") {
		t.Errorf("expected 'INFRA_' in output, got: %s", out)
	}
}
