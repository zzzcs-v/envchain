package watch

import (
	"os"
	"testing"
	"time"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envchain-watch-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New([]string{}, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New([]string{"/nonexistent/path.yaml"}, time.Second)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_ValidFile(t *testing.T) {
	p := writeTmp(t, "version: 1")
	w, err := New([]string{p}, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.hashes[p] == "" {
		t.Error("expected initial hash to be set")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	p := writeTmp(t, "version: 1")
	w, err := New([]string{p}, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	w.Start()
	defer w.Stop()

	// Overwrite the file to trigger a change.
	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("version: 2"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p {
			t.Errorf("got path %q, want %q", ev.Path, p)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoSpuriousEvents(t *testing.T) {
	p := writeTmp(t, "version: 1")
	w, err := New([]string{p}, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	w.Start()
	defer w.Stop()

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	case <-time.After(150 * time.Millisecond):
		// good — no spurious events
	}
}
