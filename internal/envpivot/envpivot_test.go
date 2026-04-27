package envpivot_test

import (
	"testing"

	"github.com/example/vaultpull/internal/envpivot"
)

func TestNew_SkipsEmptyValues(t *testing.T) {
	p := envpivot.New(map[string]string{"A": "", "B": "val"})
	if keys := p.KeysForValue(""); keys != nil {
		t.Errorf("expected nil for empty value, got %v", keys)
	}
}

func TestKeysForValue_ReturnsMatchingKeys(t *testing.T) {
	p := envpivot.New(map[string]string{"A": "shared", "B": "shared", "C": "other"})
	keys := p.KeysForValue("shared")
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestKeysForValue_ReturnsNilForMissingValue(t *testing.T) {
	p := envpivot.New(map[string]string{"A": "x"})
	if keys := p.KeysForValue("missing"); keys != nil {
		t.Errorf("expected nil, got %v", keys)
	}
}

func TestDuplicates_ReturnsDuplicateValues(t *testing.T) {
	p := envpivot.New(map[string]string{
		"DB_PASS": "secret",
		"API_KEY": "secret",
		"HOST":    "localhost",
	})
	dups := p.Duplicates()
	if len(dups) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(dups))
	}
	keys, ok := dups["secret"]
	if !ok {
		t.Fatal("expected 'secret' in duplicates")
	}
	if len(keys) != 2 || keys[0] != "API_KEY" || keys[1] != "DB_PASS" {
		t.Errorf("unexpected duplicate keys: %v", keys)
	}
}

func TestDuplicates_EmptyWhenNoDuplicates(t *testing.T) {
	p := envpivot.New(map[string]string{"A": "1", "B": "2"})
	if len(p.Duplicates()) != 0 {
		t.Error("expected no duplicates")
	}
}

func TestUniqueValues_SortedAlphabetically(t *testing.T) {
	p := envpivot.New(map[string]string{"C": "zebra", "A": "apple", "B": "mango"})
	vals := p.UniqueValues()
	if len(vals) != 3 || vals[0] != "apple" || vals[1] != "mango" || vals[2] != "zebra" {
		t.Errorf("unexpected order: %v", vals)
	}
}

func TestUniqueValues_EmptyMap(t *testing.T) {
	p := envpivot.New(map[string]string{})
	if len(p.UniqueValues()) != 0 {
		t.Error("expected empty unique values")
	}
}
