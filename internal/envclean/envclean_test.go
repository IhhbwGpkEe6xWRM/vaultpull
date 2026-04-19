package envclean_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envclean"
)

func TestClean_NoStaleKeys(t *testing.T) {
	c := envclean.New()
	local := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"A": "1", "B": "2"}
	res := c.Clean(local, incoming)
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removals, got %v", res.Removed)
	}
	if len(res.Kept) != 2 {
		t.Fatalf("expected 2 kept keys, got %d", len(res.Kept))
	}
}

func TestClean_RemovesStaleKey(t *testing.T) {
	c := envclean.New()
	local := map[string]string{"A": "1", "STALE": "old"}
	incoming := map[string]string{"A": "1"}
	res := c.Clean(local, incoming)
	if len(res.Removed) != 1 || res.Removed[0] != "STALE" {
		t.Fatalf("expected [STALE], got %v", res.Removed)
	}
	if _, ok := res.Kept["STALE"]; ok {
		t.Fatal("stale key should not be in Kept")
	}
}

func TestClean_RemovedKeysSorted(t *testing.T) {
	c := envclean.New()
	local := map[string]string{"Z": "z", "A": "a", "M": "m"}
	incoming := map[string]string{}
	res := c.Clean(local, incoming)
	if len(res.Removed) != 3 {
		t.Fatalf("expected 3 removed, got %d", len(res.Removed))
	}
	if res.Removed[0] != "A" || res.Removed[1] != "M" || res.Removed[2] != "Z" {
		t.Fatalf("unexpected order: %v", res.Removed)
	}
}

func TestClean_DryRun_KeptContainsStale(t *testing.T) {
	c := envclean.New(envclean.WithDryRun())
	local := map[string]string{"A": "1", "STALE": "old"}
	incoming := map[string]string{"A": "1"}
	res := c.Clean(local, incoming)
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 reported removal, got %d", len(res.Removed))
	}
	if _, ok := res.Kept["STALE"]; !ok {
		t.Fatal("dry-run: stale key should still be in Kept")
	}
}

func TestClean_EmptyLocal(t *testing.T) {
	c := envclean.New()
	res := c.Clean(map[string]string{}, map[string]string{"A": "1"})
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removals, got %v", res.Removed)
	}
}
