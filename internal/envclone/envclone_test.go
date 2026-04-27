package envclone_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envclone"
)

func TestClone_NilSource_ReturnsError(t *testing.T) {
	c := envclone.New()
	_, err := c.Clone(nil)
	if err == nil {
		t.Fatal("expected error for nil source, got nil")
	}
}

func TestClone_EmptyMap_ReturnsEmptyMap(t *testing.T) {
	c := envclone.New()
	out, err := c.Clone(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	c := envclone.New()
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(out))
	}
	for k, v := range src {
		if out[k] != v {
			t.Errorf("key %q: want %q, got %q", k, v, out[k])
		}
	}
}

func TestClone_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"X": "original"}
	c := envclone.New(envclone.WithValueHook(func(k, v string) (string, error) {
		return strings.ToUpper(v), nil
	}))
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["X"] != "original" {
		t.Errorf("source mutated: got %q", src["X"])
	}
	if out["X"] != "ORIGINAL" {
		t.Errorf("hook not applied: got %q", out["X"])
	}
}

func TestClone_WithKeyFilter_ExcludesKeys(t *testing.T) {
	src := map[string]string{"KEEP": "yes", "SKIP": "no"}
	c := envclone.New(envclone.WithKeyFilter(func(k string) bool {
		return k == "KEEP"
	}))
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["SKIP"]; ok {
		t.Error("SKIP should have been filtered out")
	}
	if out["KEEP"] != "yes" {
		t.Errorf("KEEP missing from clone")
	}
}

func TestClone_ValueHookError_AbortsClone(t *testing.T) {
	src := map[string]string{"BAD": "value"}
	hookErr := errors.New("hook failure")
	c := envclone.New(envclone.WithValueHook(func(k, v string) (string, error) {
		return "", hookErr
	}))
	out, err := c.Clone(src)
	if out != nil {
		t.Error("expected nil map on error")
	}
	if !errors.Is(err, hookErr) {
		t.Errorf("expected hook error, got %v", err)
	}
}
