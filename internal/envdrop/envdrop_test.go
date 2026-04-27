package envdrop_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envdrop"
)

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envdrop.New([]string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestNew_ValidPatterns_NoError(t *testing.T) {
	_, err := envdrop.New([]string{"SECRET_*", "INTERNAL_*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_RemovesMatchingKeys(t *testing.T) {
	d, _ := envdrop.New([]string{"SECRET_*"})
	input := map[string]string{
		"SECRET_KEY": "abc",
		"API_URL":    "https://example.com",
	}
	out := d.Apply(input)
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be dropped")
	}
	if out["API_URL"] != "https://example.com" {
		t.Error("expected API_URL to be preserved")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	d, _ := envdrop.New([]string{"DROP_*"})
	input := map[string]string{"DROP_ME": "val", "KEEP": "ok"}
	_ = d.Apply(input)
	if _, ok := input["DROP_ME"]; !ok {
		t.Error("Apply must not mutate the source map")
	}
}

func TestApply_NoPatterns_ReturnsAll(t *testing.T) {
	d, _ := envdrop.New(nil)
	input := map[string]string{"A": "1", "B": "2"}
	out := d.Apply(input)
	if len(out) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(out))
	}
}

func TestApply_ExactMatch(t *testing.T) {
	d, _ := envdrop.New([]string{"EXACT_KEY"})
	input := map[string]string{"EXACT_KEY": "v", "OTHER": "v2"}
	out := d.Apply(input)
	if _, ok := out["EXACT_KEY"]; ok {
		t.Error("expected EXACT_KEY to be dropped")
	}
}

func TestDropped_ReturnsSortedKeys(t *testing.T) {
	d, _ := envdrop.New([]string{"*_SECRET"})
	input := map[string]string{
		"DB_SECRET":  "x",
		"API_SECRET": "y",
		"KEEP":       "z",
	}
	keys := d.Dropped(input)
	if len(keys) != 2 {
		t.Fatalf("expected 2 dropped keys, got %d", len(keys))
	}
	if keys[0] != "API_SECRET" || keys[1] != "DB_SECRET" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestDropped_EmptyWhenNoMatch(t *testing.T) {
	d, _ := envdrop.New([]string{"NOPE_*"})
	input := map[string]string{"A": "1"}
	keys := d.Dropped(input)
	if len(keys) != 0 {
		t.Errorf("expected no dropped keys, got %v", keys)
	}
}
