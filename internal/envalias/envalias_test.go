package envalias_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envalias"
)

func TestNew_ParsesPairs(t *testing.T) {
	m := envalias.New([]string{"DB_PASS=DATABASE_PASSWORD", "API_KEY=SERVICE_API_KEY"})
	if len(m.Pairs()) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(m.Pairs()))
	}
}

func TestNew_IgnoresMalformed(t *testing.T) {
	m := envalias.New([]string{"NODASH", "=empty_from", "empty_to=", "GOOD=OK"})
	if len(m.Pairs()) != 1 {
		t.Fatalf("expected 1 valid pair, got %d", len(m.Pairs()))
	}
}

func TestApply_RenamesKey(t *testing.T) {
	m := envalias.New([]string{"OLD_KEY=NEW_KEY"})
	result := m.Apply(map[string]string{"OLD_KEY": "secret"})
	if result["NEW_KEY"] != "secret" {
		t.Errorf("expected NEW_KEY=secret, got %q", result["NEW_KEY"])
	}
	if _, ok := result["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
}

func TestApply_SkipsMissingFrom(t *testing.T) {
	m := envalias.New([]string{"MISSING=TARGET"})
	result := m.Apply(map[string]string{"OTHER": "val"})
	if _, ok := result["TARGET"]; ok {
		t.Error("TARGET should not be present")
	}
	if result["OTHER"] != "val" {
		t.Error("OTHER should be preserved")
	}
}

func TestApply_DoesNotOverwriteExistingTo(t *testing.T) {
	m := envalias.New([]string{"SRC=DEST"})
	result := m.Apply(map[string]string{"SRC": "new", "DEST": "original"})
	if result["DEST"] != "original" {
		t.Errorf("DEST should not be overwritten, got %q", result["DEST"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	m := envalias.New([]string{"A=B"})
	input := map[string]string{"A": "val"}
	m.Apply(input)
	if _, ok := input["A"]; !ok {
		t.Error("original input map should not be mutated")
	}
}

func TestApply_EmptyAliases_ReturnsClone(t *testing.T) {
	m := envalias.New(nil)
	input := map[string]string{"X": "1", "Y": "2"}
	result := m.Apply(input)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}
