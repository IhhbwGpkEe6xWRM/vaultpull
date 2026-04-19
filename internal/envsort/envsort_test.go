package envsort_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envsort"
)

func TestKeys_AlphaOrder(t *testing.T) {
	s := envsort.New()
	m := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	keys := s.Keys(m)
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Fatalf"index %d: got %q want %q", i, k, want[i])
		}
	}
}

func TestKeys_PriorityFirst(t *testing.T) {
	s := envsort.New(envsort.WithPriority("DB_URL", "APP_ENV"))
	m := map[string]string{"ZEBRA": "z", "DB_URL": "postgres", "APP_ENV": "prod", "ALPHA": "a"}
	keys := s.Keys(m)
	if keys[0] != "DB_URL" {
		t.Fatalf("expected DB_URL first, got %q", keys[0])
	}
	if keys[1] != "APP_ENV" {
		t.Fatalf("expected APP_ENV second, got %q", keys[1])
	}
}

func TestKeys_PriorityKeyNotInMap(t *testing.T) {
	s := envsort.New(envsort.WithPriority("MISSING"))
	m := map[string]string{"B": "1", "A": "2"}
	keys := s.Keys(m)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "B" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestApply_FormatsKeyValue(t *testing.T) {
	s := envsort.New()
	m := map[string]string{"FOO": "bar"}
	out := s.Apply(m)
	if len(out) != 1 || out[0] != "FOO=bar" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	s := envsort.New()
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %v", out)
	}
}

func TestKeys_StableWithNoPriority(t *testing.T) {
	s := envsort.New()
	m := map[string]string{"C": "3", "A": "1", "B": "2"}
	for i := 0; i < 10; i++ {
		keys := s.Keys(m)
		if keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
			t.Fatalf("unstable order on iteration %d: %v", i, keys)
		}
	}
}
