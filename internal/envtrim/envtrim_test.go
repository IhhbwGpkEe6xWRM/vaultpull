package envtrim_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envtrim"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	src := map[string]string{"A": "", "B": "  ", "C": "val"}
	out := envtrim.New().Apply(src)
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestApply_TrimEmpty_RemovesEmptyValues(t *testing.T) {
	src := map[string]string{"A": "", "B": "hello", "C": ""}
	out := envtrim.New(envtrim.WithTrimEmpty()).Apply(src)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["B"] != "hello" {
		t.Errorf("expected B=hello, got %q", out["B"])
	}
}

func TestApply_TrimWhitespace_RemovesBlankValues(t *testing.T) {
	src := map[string]string{"A": "  ", "B": "\t", "C": "real"}
	out := envtrim.New(envtrim.WithTrimWhitespace()).Apply(src)
	if _, ok := out["A"]; ok {
		t.Error("expected A to be removed")
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if out["C"] != "real" {
		t.Errorf("expected C=real, got %q", out["C"])
	}
}

func TestApply_TrimFunc_CustomPredicate(t *testing.T) {
	src := map[string]string{"A": "REDACTED", "B": "value", "C": "REDACTED"}
	out := envtrim.New(envtrim.WithTrimFunc(func(v string) bool {
		return v == "REDACTED"
	})).Apply(src)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"A": "", "B": "keep"}
	envtrim.New(envtrim.WithTrimEmpty()).Apply(src)
	if len(src) != 2 {
		t.Error("source map was mutated")
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	src := map[string]string{"A": "", "B": "  ", "C": "skip", "D": "keep"}
	out := envtrim.New(
		envtrim.WithTrimEmpty(),
		envtrim.WithTrimWhitespace(),
		envtrim.WithTrimFunc(func(v string) bool { return strings.ToLower(v) == "skip" }),
	).Apply(src)
	if len(out) != 1 || out["D"] != "keep" {
		t.Errorf("expected only D=keep, got %v", out)
	}
}

func TestRemoved_ReturnsRemovedKeys(t *testing.T) {
	src := map[string]string{"X": "", "Y": "present", "Z": ""}
	removed := envtrim.New(envtrim.WithTrimEmpty()).Removed(src)
	sort.Strings(removed)
	if len(removed) != 2 || removed[0] != "X" || removed[1] != "Z" {
		t.Errorf("unexpected removed keys: %v", removed)
	}
}

func TestRemoved_EmptyMap_ReturnsNil(t *testing.T) {
	removed := envtrim.New(envtrim.WithTrimEmpty()).Removed(map[string]string{})
	if len(removed) != 0 {
		t.Errorf("expected no removed keys, got %v", removed)
	}
}
