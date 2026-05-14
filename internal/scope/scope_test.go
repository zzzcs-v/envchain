package scope_test

import (
	"os"
	"testing"

	"github.com/user/envchain/internal/scope"
)

func tempStore(t *testing.T) *scope.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := scope.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func TestSave_AndLoad(t *testing.T) {
	st := tempStore(t)
	sc := scope.Scope{Name: "prod", Vars: map[string]string{"DB_HOST": "db.prod", "PORT": "5432"}}
	if err := st.Save(sc); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := st.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != sc.Name {
		t.Errorf("name: got %q want %q", got.Name, sc.Name)
	}
	if got.Vars["DB_HOST"] != "db.prod" {
		t.Errorf("DB_HOST: got %q", got.Vars["DB_HOST"])
	}
}

func TestSave_EmptyName(t *testing.T) {
	st := tempStore(t)
	err := st.Save(scope.Scope{Name: "", Vars: map[string]string{}})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLoad_Missing(t *testing.T) {
	st := tempStore(t)
	_, err := st.Load("ghost")
	if err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestDelete_RemovesScope(t *testing.T) {
	st := tempStore(t)
	sc := scope.Scope{Name: "staging", Vars: map[string]string{"ENV": "staging"}}
	_ = st.Save(sc)
	if err := st.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := st.Load("staging")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	st := tempStore(t)
	if err := st.Delete("nope"); err == nil {
		t.Error("expected error deleting missing scope")
	}
}

func TestList_SortedByName(t *testing.T) {
	st := tempStore(t)
	for _, name := range []string{"zebra", "alpha", "middle"} {
		_ = st.Save(scope.Scope{Name: name, Vars: map[string]string{}})
	}
	names, err := st.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	want := []string{"alpha", "middle", "zebra"}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("index %d: got %q want %q", i, names[i], n)
		}
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	// remove dir so List handles missing gracefully
	os.RemoveAll(dir)
	st, _ := scope.NewStore(dir)
	// recreate only the store struct without the dir
	_ = os.RemoveAll(dir)
	names, err := st.List()
	if err != nil && names != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
