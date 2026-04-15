package diff_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/diff"
)

func TestCompare_AllUnchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	next := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := diff.Compare(old, next)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	a, rem, m := r.Summary()
	if a != 0 || rem != 0 || m != 0 {
		t.Fatalf("unexpected summary: added=%d removed=%d modified=%d", a, rem, m)
	}
}

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	next := map[string]string{"NEW_KEY": "value"}
	r := diff.Compare(old, next)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	a, _, _ := r.Summary()
	if a != 1 {
		t.Fatalf("expected 1 added, got %d", a)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"OLD_KEY": "value"}
	next := map[string]string{}
	r := diff.Compare(old, next)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	_, rem, _ := r.Summary()
	if rem != 1 {
		t.Fatalf("expected 1 removed, got %d", rem)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := map[string]string{"KEY": "old"}
	next := map[string]string{"KEY": "new"}
	r := diff.Compare(old, next)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	_, _, m := r.Summary()
	if m != 1 {
		t.Fatalf("expected 1 modified, got %d", m)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "1"}
	next := map[string]string{"Z": "2", "A": "1", "M": "3"}
	r := diff.Compare(old, next)
	keys := make([]string, len(r.Changes))
	for i, c := range r.Changes {
		keys[i] = c.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("changes not sorted: %v", keys)
		}
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	old := map[string]string{"KEEP": "same", "CHANGE": "old", "DROP": "gone"}
	next := map[string]string{"KEEP": "same", "CHANGE": "new", "ADD": "fresh"}
	r := diff.Compare(old, next)
	a, rem, m := r.Summary()
	if a != 1 || rem != 1 || m != 1 {
		t.Fatalf("expected 1/1/1 got added=%d removed=%d modified=%d", a, rem, m)
	}
}
