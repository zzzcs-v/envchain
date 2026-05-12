package pin_test

import (
	"os"
	"testing"

	"github.com/nicholasgasior/envchain/internal/pin"
)

func tempStore(t *testing.T) *pin.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "pin-test-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := pin.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestSave_AndLoad(t *testing.T) {
	s := tempStore(t)
	p := pin.Pin{
		Name:    "mypin",
		Context: "staging",
		Vars:    map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("mypin")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Context != "staging" {
		t.Errorf("context: want staging, got %s", got.Context)
	}
	if got.Vars["FOO"] != "bar" {
		t.Errorf("FOO: want bar, got %s", got.Vars["FOO"])
	}
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestSave_EmptyName(t *testing.T) {
	s := tempStore(t)
	err := s.Save(pin.Pin{Name: "", Context: "prod"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestLoad_Missing(t *testing.T) {
	s := tempStore(t)
	_, err := s.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing pin")
	}
}

func TestDelete_RemovesPin(t *testing.T) {
	s := tempStore(t)
	p := pin.Pin{Name: "todelete", Context: "dev", Vars: map[string]string{}}
	_ = s.Save(p)
	if err := s.Delete("todelete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Load("todelete")
	if err == nil {
		t.Fatal("expected pin to be gone after delete")
	}
}

func TestList_ReturnsSaved(t *testing.T) {
	s := tempStore(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = s.Save(pin.Pin{Name: name, Context: "dev", Vars: map[string]string{"K": "v"}})
	}
	pins, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(pins) != 3 {
		t.Errorf("want 3 pins, got %d", len(pins))
	}
}
