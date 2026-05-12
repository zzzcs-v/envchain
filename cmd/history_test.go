package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/history"
)

func setupHistoryStore(t *testing.T) (string, *history.Store) {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-hist-cmd-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return dir, s
}

func TestRunHistory_Empty(t *testing.T) {
	dir, _ := setupHistoryStore(t)
	historyDir = dir
	var buf bytes.Buffer
	historyCmd.SetOut(&buf)
	if err := runHistory(historyCmd, nil); err != nil {
		t.Fatalf("runHistory: %v", err)
	}
	if !strings.Contains(buf.String(), "no history entries") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestRunHistory_ShowsEntries(t *testing.T) {
	dir, s := setupHistoryStore(t)
	historyDir = dir
	_ = s.Record("dev", "dotenv", map[string]string{"A": "1", "B": "2"})
	var buf bytes.Buffer
	historyCmd.SetOut(&buf)
	if err := runHistory(historyCmd, nil); err != nil {
		t.Fatalf("runHistory: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dev") {
		t.Errorf("expected context dev in output: %q", out)
	}
	if !strings.Contains(out, "dotenv") {
		t.Errorf("expected format dotenv in output: %q", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected var count 2 in output: %q", out)
	}
}

func TestRunHistoryClear(t *testing.T) {
	dir, s := setupHistoryStore(t)
	historyDir = dir
	_ = s.Record("prod", "json", nil)
	var buf bytes.Buffer
	historyClearCmd.SetOut(&buf)
	if err := runHistoryClear(historyClearCmd, nil); err != nil {
		t.Fatalf("runHistoryClear: %v", err)
	}
	if !strings.Contains(buf.String(), "cleared") {
		t.Errorf("expected cleared message, got: %q", buf.String())
	}
	entries, _ := s.List()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}
