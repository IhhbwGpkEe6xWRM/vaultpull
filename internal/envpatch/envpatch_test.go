package envpatch_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envpatch"
)

func TestNew_ValidEntries(t *testing.T) {
	_, err := envpatch.New([]envpatch.Entry{
		{Key: "FOO", Value: "bar", Op: envpatch.OpSet},
		{Key: "BAZ", Op: envpatch.OpDelete},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envpatch.New([]envpatch.Entry{{Key: "", Op: envpatch.OpSet}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_UnknownOp_ReturnsError(t *testing.T) {
	_, err := envpatch.New([]envpatch.Entry{{Key: "X", Op: "upsert"}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestApply_SetsNewKey(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Entry{{Key: "NEW", Value: "val", Op: envpatch.OpSet}})
	out := p.Apply(map[string]string{"EXISTING": "x"})
	if out["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %q", out["NEW"])
	}
	if out["EXISTING"] != "x" {
		t.Errorf("expected EXISTING=x, got %q", out["EXISTING"])
	}
}

func TestApply_OverwritesExistingKey(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Entry{{Key: "K", Value: "new", Op: envpatch.OpSet}})
	out := p.Apply(map[string]string{"K": "old"})
	if out["K"] != "new" {
		t.Errorf("expected K=new, got %q", out["K"])
	}
}

func TestApply_DeletesKey(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Entry{{Key: "GONE", Op: envpatch.OpDelete}})
	out := p.Apply(map[string]string{"GONE": "bye", "KEEP": "yes"})
	if _, ok := out["GONE"]; ok {
		t.Error("expected GONE to be deleted")
	}
	if out["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %q", out["KEEP"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	p, _ := envpatch.New([]envpatch.Entry{{Key: "A", Value: "2", Op: envpatch.OpSet}})
	p.Apply(base)
	if base["A"] != "1" {
		t.Error("Apply must not mutate the base map")
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	entries := envpatch.Diff(map[string]string{}, map[string]string{"X": "1"})
	if len(entries) != 1 || entries[0].Key != "X" || entries[0].Op != envpatch.OpSet {
		t.Errorf("unexpected diff: %+v", entries)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	entries := envpatch.Diff(map[string]string{"X": "1"}, map[string]string{})
	if len(entries) != 1 || entries[0].Key != "X" || entries[0].Op != envpatch.OpDelete {
		t.Errorf("unexpected diff: %+v", entries)
	}
}

func TestDiff_UnchangedKeys_NotIncluded(t *testing.T) {
	entries := envpatch.Diff(
		map[string]string{"A": "1", "B": "2"},
		map[string]string{"A": "1", "B": "2"},
	)
	if len(entries) != 0 {
		t.Errorf("expected empty diff, got %+v", entries)
	}
}
