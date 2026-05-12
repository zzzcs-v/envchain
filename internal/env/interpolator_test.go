package env

import (
	"testing"
)

func makeInterp(vars map[string]string) *Interpolator {
	i := New(vars)
	// override OS lookup so tests are hermetic
	i.osLookup = func(key string) (string, bool) {
		if key == "OS_VAR" {
			return "from_os", true
		}
		return "", false
	}
	return i
}

func TestResolve_NoPlaceholders(t *testing.T) {
	i := makeInterp(nil)
	out, err := i.Resolve("hello world")
	if err != nil || out != "hello world" {
		t.Fatalf("expected 'hello world', got %q err=%v", out, err)
	}
}

func TestResolve_FromVarsMap(t *testing.T) {
	i := makeInterp(map[string]string{"HOST": "localhost"})
	out, err := i.Resolve("http://${HOST}:8080")
	if err != nil {
		t.Fatal(err)
	}
	if out != "http://localhost:8080" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestResolve_FromOS(t *testing.T) {
	i := makeInterp(nil)
	out, err := i.Resolve("val=${OS_VAR}")
	if err != nil {
		t.Fatal(err)
	}
	if out != "val=from_os" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestResolve_InlineDefault(t *testing.T) {
	i := makeInterp(nil)
	out, err := i.Resolve("${MISSING:-fallback}")
	if err != nil {
		t.Fatal(err)
	}
	if out != "fallback" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestResolve_Unresolved_ReturnsError(t *testing.T) {
	i := makeInterp(nil)
	_, err := i.Resolve("${NOPE}")
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestResolveAll_Success(t *testing.T) {
	i := makeInterp(map[string]string{"BASE": "https://example.com"})
	in := map[string]string{
		"API_URL": "${BASE}/api",
		"STATIC":  "plain",
	}
	out, err := i.ResolveAll(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("unexpected API_URL: %q", out["API_URL"])
	}
	if out["STATIC"] != "plain" {
		t.Errorf("unexpected STATIC: %q", out["STATIC"])
	}
}

func TestResolveAll_PropagatesError(t *testing.T) {
	i := makeInterp(nil)
	_, err := i.ResolveAll(map[string]string{"KEY": "${GHOST}"})
	if err == nil {
		t.Fatal("expected error")
	}
}
