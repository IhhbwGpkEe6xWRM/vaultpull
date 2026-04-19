package envmerge

import (
	"testing"
)

func TestMerge_VaultWins_NoConflict(t *testing.T) {
	local := map[string]string{"A": "1"}
	vault := map[string]string{"B": "2"}
	r, err := Merge(local, vault, VaultWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Secrets["A"] != "1" || r.Secrets["B"] != "2" {
		t.Errorf("unexpected secrets: %v", r.Secrets)
	}
	if len(r.Overridden) != 0 {
		t.Errorf("expected no overrides, got %v", r.Overridden)
	}
}

func TestMerge_VaultWins_OverridesLocal(t *testing.T) {
	local := map[string]string{"KEY": "local"}
	vault := map[string]string{"KEY": "vault"}
	r, err := Merge(local, vault, VaultWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Secrets["KEY"] != "vault" {
		t.Errorf("expected vault value, got %q", r.Secrets["KEY"])
	}
	if len(r.Overridden) != 1 || r.Overridden[0] != "KEY" {
		t.Errorf("expected KEY in overridden, got %v", r.Overridden)
	}
}

func TestMerge_LocalWins_KeepsLocal(t *testing.T) {
	local := map[string]string{"KEY": "local"}
	vault := map[string]string{"KEY": "vault"}
	r, err := Merge(local, vault, LocalWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Secrets["KEY"] != "local" {
		t.Errorf("expected local value, got %q", r.Secrets["KEY"])
	}
}

func TestMerge_ErrorOnConflict_ReturnsError(t *testing.T) {
	local := map[string]string{"X": "a", "Y": "same"}
	vault := map[string]string{"X": "b", "Y": "same"}
	_, err := Merge(local, vault, ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "envmerge: conflicts on keys: X" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMerge_SameValue_NoConflict(t *testing.T) {
	local := map[string]string{"K": "v"}
	vault := map[string]string{"K": "v"}
	r, err := Merge(local, vault, ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Overridden) != 0 {
		t.Errorf("expected no overrides for identical values")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	r, err := Merge(nil, nil, VaultWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Secrets) != 0 {
		t.Errorf("expected empty result")
	}
}
