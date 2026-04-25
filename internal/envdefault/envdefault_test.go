package envdefault_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envdefault"
)

func newApplier(t *testing.T, pairs []string, opts ...envdefault.Option) *envdefault.Applier {
	t.Helper()
	a, err := envdefault.New(pairs, opts...)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	return a
}

func TestNew_ValidPairs(t *testing.T) {
	a := newApplier(t, []string{"FOO=bar", "BAZ="})
	keys := a.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestNew_MalformedPair(t *testing.T) {
	_, err := envdefault.New([]string{"NOEQUALSSIGN"})
	if err == nil {
		t.Fatal("expected error for malformed pair, got nil")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := envdefault.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestApply_FillsMissingKey(t *testing.T) {
	a := newApplier(t, []string{"FOO=default"})
	out := a.Apply(map[string]string{})
	if out["FOO"] != "default" {
		t.Errorf("expected 'default', got %q", out["FOO"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	a := newApplier(t, []string{"FOO=default"})
	out := a.Apply(map[string]string{"FOO": "existing"})
	if out["FOO"] != "existing" {
		t.Errorf("expected 'existing', got %q", out["FOO"])
	}
}

func TestApply_FillsEmptyValue(t *testing.T) {
	a := newApplier(t, []string{"FOO=default"})
	out := a.Apply(map[string]string{"FOO": ""})
	if out["FOO"] != "default" {
		t.Errorf("expected 'default', got %q", out["FOO"])
	}
}

func TestApply_WithOverwrite_ReplacesExisting(t *testing.T) {
	a := newApplier(t, []string{"FOO=default"}, envdefault.WithOverwrite())
	out := a.Apply(map[string]string{"FOO": "existing"})
	if out["FOO"] != "default" {
		t.Errorf("expected 'default', got %q", out["FOO"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	a := newApplier(t, []string{"FOO=default"})
	src := map[string]string{"BAR": "baz"}
	a.Apply(src)
	if _, ok := src["FOO"]; ok {
		t.Error("Apply mutated the input map")
	}
}

func TestApply_EmptyDefaults_ReturnsClone(t *testing.T) {
	a := newApplier(t, []string{})
	src := map[string]string{"A": "1"}
	out := a.Apply(src)
	if out["A"] != "1" {
		t.Errorf("expected '1', got %q", out["A"])
	}
}
