package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	from := map[string]string{"FOO": "bar", "BAZ": "qux"}
	to := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(from, to)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_Added(t *testing.T) {
	from := map[string]string{}
	to := map[string]string{"NEW_KEY": "value"}
	r := Compare(from, to)
	if len(r.Changes) != 1 || r.Changes[0].Type != Added {
		t.Errorf("expected one Added change, got %+v", r.Changes)
	}
	if r.Changes[0].NewValue != "value" {
		t.Errorf("unexpected new value: %s", r.Changes[0].NewValue)
	}
}

func TestCompare_Removed(t *testing.T) {
	from := map[string]string{"OLD_KEY": "gone"}
	to := map[string]string{}
	r := Compare(from, to)
	if len(r.Changes) != 1 || r.Changes[0].Type != Removed {
		t.Errorf("expected one Removed change, got %+v", r.Changes)
	}
}

func TestCompare_Modified(t *testing.T) {
	from := map[string]string{"HOST": "localhost"}
	to := map[string]string{"HOST": "prod.example.com"}
	r := Compare(from, to)
	if len(r.Changes) != 1 || r.Changes[0].Type != Modified {
		t.Errorf("expected one Modified change, got %+v", r.Changes)
	}
	if r.Changes[0].OldValue != "localhost" || r.Changes[0].NewValue != "prod.example.com" {
		t.Errorf("unexpected values in change: %+v", r.Changes[0])
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	from := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	to := map[string]string{"Z_KEY": "1", "A_KEY": "99"}
	r := Compare(from, to)
	if r.Changes[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %s", r.Changes[0].Key)
	}
}

func TestResult_Summary(t *testing.T) {
	from := map[string]string{"A": "1", "B": "old"}
	to := map[string]string{"B": "new", "C": "3"}
	r := Compare(from, to)
	summary := r.Summary()
	if summary != "+1 added, -1 removed, ~1 modified" {
		t.Errorf("unexpected summary: %s", summary)
	}
}
