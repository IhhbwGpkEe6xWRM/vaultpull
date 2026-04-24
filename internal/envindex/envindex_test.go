package envindex_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envindex"
)

func TestNew_SkipsEmptyValues(t *testing.T) {
	idx := envindex.New(map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
	})
	if idx.Len() != 1 {
		t.Fatalf("expected 1 distinct value, got %d", idx.Len())
	}
}

func TestKeysForValue_ReturnsMatchingKeys(t *testing.T) {
	idx := envindex.New(map[string]string{
		"KEY_A": "shared",
		"KEY_B": "shared",
		"KEY_C": "unique",
	})
	keys := idx.KeysForValue("shared")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "KEY_A" || keys[1] != "KEY_B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestKeysForValue_ReturnsNilForMissingValue(t *testing.T) {
	idx := envindex.New(map[string]string{"KEY": "val"})
	if got := idx.KeysForValue("missing"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestDuplicates_ReturnsDuplicateValues(t *testing.T) {
	idx := envindex.New(map[string]string{
		"A": "dup",
		"B": "dup",
		"C": "solo",
	})
	dups := idx.Duplicates()
	if len(dups) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(dups))
	}
	keys, ok := dups["dup"]
	if !ok {
		t.Fatal("expected 'dup' in duplicates map")
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys for 'dup', got %d", len(keys))
	}
}

func TestDuplicates_EmptyWhenNoDuplicates(t *testing.T) {
	idx := envindex.New(map[string]string{
		"A": "one",
		"B": "two",
	})
	if len(idx.Duplicates()) != 0 {
		t.Error("expected no duplicates")
	}
}

func TestHasDuplicates_True(t *testing.T) {
	idx := envindex.New(map[string]string{
		"X": "same",
		"Y": "same",
	})
	if !idx.HasDuplicates() {
		t.Error("expected HasDuplicates to return true")
	}
}

func TestHasDuplicates_False(t *testing.T) {
	idx := envindex.New(map[string]string{
		"X": "a",
		"Y": "b",
	})
	if idx.HasDuplicates() {
		t.Error("expected HasDuplicates to return false")
	}
}

func TestNew_EmptyMap(t *testing.T) {
	idx := envindex.New(map[string]string{})
	if idx.Len() != 0 {
		t.Errorf("expected length 0, got %d", idx.Len())
	}
	if idx.HasDuplicates() {
		t.Error("empty index should have no duplicates")
	}
}
