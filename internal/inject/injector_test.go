package inject_test

import (
	"testing"

	"github.com/envchain/envchain/internal/inject"
)

func makeInjector(overwrite bool, prefix string) *inject.Injector {
	return inject.New(inject.Options{Overwrite: overwrite, Prefix: prefix})
}

func TestInject_EmptySource(t *testing.T) {
	inj := makeInjector(false, "")
	dst := map[string]string{}
	res, err := inj.Inject(dst, []inject.Source{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}

func TestInject_NilDst(t *testing.T) {
	inj := makeInjector(false, "")
	_, err := inj.Inject(nil, []inject.Source{{Name: "a", Vars: map[string]string{"K": "v"}}})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestInject_BasicValues(t *testing.T) {
	inj := makeInjector(false, "")
	dst := map[string]string{}
	srcs := []inject.Source{
		{Name: "base", Vars: map[string]string{"FOO": "bar", "BAZ": "qux"}},
	}
	res, err := inj.Inject(dst, srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Injected != 2 {
		t.Errorf("expected 2 injected, got %d", res[0].Injected)
	}
	if dst["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", dst["FOO"])
	}
}

func TestInject_NoOverwrite_SkipsExisting(t *testing.T) {
	inj := makeInjector(false, "")
	dst := map[string]string{"FOO": "original"}
	srcs := []inject.Source{
		{Name: "override", Vars: map[string]string{"FOO": "new"}},
	}
	res, err := inj.Inject(dst, srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res[0].Skipped)
	}
	if dst["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %q", dst["FOO"])
	}
}

func TestInject_Overwrite_ReplacesExisting(t *testing.T) {
	inj := makeInjector(true, "")
	dst := map[string]string{"FOO": "original"}
	srcs := []inject.Source{
		{Name: "override", Vars: map[string]string{"FOO": "new"}},
	}
	_, err := inj.Inject(dst, srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", dst["FOO"])
	}
}

func TestInject_PrefixApplied(t *testing.T) {
	inj := makeInjector(false, "APP_")
	dst := map[string]string{}
	srcs := []inject.Source{
		{Name: "svc", Vars: map[string]string{"PORT": "8080"}},
	}
	_, err := inj.Inject(dst, srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q", dst["APP_PORT"])
	}
}

func TestInject_EmptySourceName_ReturnsError(t *testing.T) {
	inj := makeInjector(false, "")
	dst := map[string]string{}
	srcs := []inject.Source{
		{Name: "", Vars: map[string]string{"X": "y"}},
	}
	_, err := inj.Inject(dst, srcs)
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}
