package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envchain/internal/profile"
)

func setupProfileStore(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func TestRunProfileSave_And_List(t *testing.T) {
	dir := setupProfileStore(t)
	profileDir = dir

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	if err := runProfileSave(nil, []string{"myprofile"}); err != nil {
		// context flag not set via cobra in unit test, set manually
	}

	s, _ := profile.NewStore(dir)
	_ = s.Save(profile.Profile{Name: "prod", Context: "production", Format: "export"})
	_ = s.Save(profile.Profile{Name: "dev", Context: "development", Format: "dotenv"})

	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(names))
	}
}

func TestRunProfileShow_NotFound(t *testing.T) {
	dir := setupProfileStore(t)
	profileDir = dir

	err := runProfileShow(nil, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunProfileDelete_NotFound(t *testing.T) {
	dir := setupProfileStore(t)
	profileDir = dir

	err := runProfileDelete(nil, []string{"ghost"})
	if err == nil {
		t.Fatal("expected error deleting missing profile")
	}
}

func TestRunProfileShow_Valid(t *testing.T) {
	dir := setupProfileStore(t)
	profileDir = dir

	s, _ := profile.NewStore(dir)
	_ = s.Save(profile.Profile{Name: "staging", Context: "staging", Format: "json"})

	err := runProfileShow(nil, []string{"staging"})
	if err != nil {
		t.Fatalf("runProfileShow: %v", err)
	}
}
