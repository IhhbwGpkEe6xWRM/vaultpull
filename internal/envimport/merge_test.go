package envimport_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpull/internal/envimport"
)

func TestMerge_VaultWins_OverridesLocal(t *testing.T) {
	local := map[string]string{"KEY": "old"}
	vault := map[string]string{"KEY": "new"}
	out, err := envimport.Merge(local, vault, envimport.VaultWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "new" {
		t.Fatalf("expected new, got %q", out["KEY"])
	}
}

func TestMerge_LocalWins_KeepsLocal(t *testing.T) {
	local := map[string]string{"KEY": "local"}
	vault := map[string]string{"KEY": "vault"}
	out, err := envimport.Merge(local, vault, envimport.LocalWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "local" {
		t.Fatalf("expected local, got %q", out["KEY"])
	}
}

func TestMerge_ErrorOnConflict_ReturnsError(t *testing.T) {
	local := map[string]string{"KEY": "a"}
	vault := map[string]string{"KEY": "b"}
	_, err := envimport.Merge(local, vault, envimport.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error")
	}
	var ce *envimport.ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected ConflictError, got %T", err)
	}
	if ce.Key != "KEY" {
		t.Fatalf("unexpected key: %q", ce.Key)
	}
}

func TestMerge_SameValue_NoConflict(t *testing.T) {
	local := map[string]string{"KEY": "same"}
	vault := map[string]string{"KEY": "same"}
	out, err := envimport.Merge(local, vault, envimport.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "same" {
		t.Fatalf("unexpected value: %q", out["KEY"])
	}
}

func TestMerge_UniqueKeys_Combined(t *testing.T) {
	local := map[string]string{"A": "1"}
	vault := map[string]string{"B": "2"}
	out, err := envimport.Merge(local, vault, envimport.VaultWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("unexpected map: %v", out)
	}
}
