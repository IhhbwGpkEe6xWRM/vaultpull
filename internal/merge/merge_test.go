package merge_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/merge"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	res, err := merge.Merge(merge.LastWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["FOO"] != "1" || res.Secrets["BAR"] != "2" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_LastWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res, err := merge.Merge(merge.LastWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Secrets["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res, err := merge.Merge(merge.FirstWins, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Secrets["KEY"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "a"}
	b := map[string]string{"KEY": "b"}
	_, err := merge.Merge(merge.ErrorOnConflict, a, b)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := merge.Merge(merge.LastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestMerge_MultipleSources_LastWins(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	c := map[string]string{"X": "3"}
	res, err := merge.Merge(merge.LastWins, a, b, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["X"] != "3" {
		t.Errorf("expected '3', got %q", res.Secrets["X"])
	}
}
