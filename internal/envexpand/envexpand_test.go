package envexpand_test

import (
	"os"
	"testing"

	"github.com/hashicorp/vaultpull/internal/envexpand"
)

func TestExpand_NoReferences(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := e.Expand(in)
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestExpand_ResolvesReference(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out := e.Expand(in)
	if want := "https://example.com/api"; out["API_URL"] != want {
		t.Fatalf("got %q, want %q", out["API_URL"], want)
	}
}

func TestExpand_SelfReferenceIsEmpty(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{"FOO": "${FOO}_suffix"}
	out := e.Expand(in)
	if out["FOO"] != "_suffix" {
		t.Fatalf("got %q, want %q", out["FOO"], "_suffix")
	}
}

func TestExpand_MissingKeyBecomesEmpty(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{"FOO": "${MISSING}_val"}
	out := e.Expand(in)
	if out["FOO"] != "_val" {
		t.Fatalf("got %q, want %q", out["FOO"], "_val")
	}
}

func TestExpand_WithOSFallback(t *testing.T) {
	os.Setenv("_TEST_OS_KEY", "fromenv")
	t.Cleanup(func() { os.Unsetenv("_TEST_OS_KEY") })

	e := envexpand.New(envexpand.WithOSFallback())
	in := map[string]string{"FOO": "${_TEST_OS_KEY}_val"}
	out := e.Expand(in)
	if out["FOO"] != "fromenv_val" {
		t.Fatalf("got %q, want %q", out["FOO"], "fromenv_val")
	}
}

func TestExpand_WithoutOSFallback_IgnoresEnv(t *testing.T) {
	os.Setenv("_TEST_OS_KEY2", "fromenv")
	t.Cleanup(func() { os.Unsetenv("_TEST_OS_KEY2") })

	e := envexpand.New()
	in := map[string]string{"FOO": "${_TEST_OS_KEY2}"}
	out := e.Expand(in)
	if out["FOO"] != "" {
		t.Fatalf("expected empty, got %q", out["FOO"])
	}
}

func TestContainsReferences_True(t *testing.T) {
	if !envexpand.ContainsReferences("${FOO}") {
		t.Fatal("expected true")
	}
}

func TestContainsReferences_False(t *testing.T) {
	if envexpand.ContainsReferences("plainvalue") {
		t.Fatal("expected false")
	}
}
