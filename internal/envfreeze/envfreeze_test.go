package envfreeze_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/envfreeze"
)

func TestNew_EmptySource(t *testing.T) {
	f := envfreeze.New(map[string]string{})
	if len(f.Keys()) != 0 {
		t.Fatalf("expected no frozen keys, got %v", f.Keys())
	}
}

func TestIsFrozen_ReturnsTrueForFrozenKey(t *testing.T) {
	f := envfreeze.New(map[string]string{"DB_PASS": "secret"})
	if !f.IsFrozen("DB_PASS") {
		t.Fatal("expected DB_PASS to be frozen")
	}
}

func TestIsFrozen_ReturnsFalseForUnknownKey(t *testing.T) {
	f := envfreeze.New(map[string]string{"DB_PASS": "secret"})
	if f.IsFrozen("OTHER") {
		t.Fatal("expected OTHER not to be frozen")
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	f := envfreeze.New(map[string]string{"Z": "1", "A": "2", "M": "3"})
	keys := f.Keys()
	want := []string{"A", "M", "Z"}
	for i, k := range want {
		if keys[i] != k {
			t.Fatalf("keys[%d] = %q, want %q", i, keys[i], k)
		}
	}
}

func TestApply_NonFrozenKeyMergedFreely(t *testing.T) {
	f := envfreeze.New(map[string]string{"FROZEN": "x"})
	out, err := f.Apply(
		map[string]string{"FROZEN": "x"},
		map[string]string{"NEW_KEY": "hello"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "hello" {
		t.Fatalf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
}

func TestApply_FrozenKeyUnchanged_OK(t *testing.T) {
	f := envfreeze.New(map[string]string{"TOKEN": "abc"})
	_, err := f.Apply(
		map[string]string{"TOKEN": "abc"},
		map[string]string{"TOKEN": "abc"},
	)
	if err != nil {
		t.Fatalf("unexpected error for unchanged frozen key: %v", err)
	}
}

func TestApply_FrozenKeyChanged_ReturnsError(t *testing.T) {
	f := envfreeze.New(map[string]string{"TOKEN": "abc"})
	_, err := f.Apply(
		map[string]string{"TOKEN": "abc"},
		map[string]string{"TOKEN": "xyz"},
	)
	if !errors.Is(err, envfreeze.ErrFrozen) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	f := envfreeze.New(map[string]string{})
	base := map[string]string{"A": "1"}
	_, _ = f.Apply(base, map[string]string{"B": "2"})
	if _, ok := base["B"]; ok {
		t.Fatal("Apply must not mutate the base map")
	}
}
