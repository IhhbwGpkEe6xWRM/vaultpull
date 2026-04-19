package envrewrite_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envrewrite"
)

func TestNew_InvalidRule_EmptyFrom(t *testing.T) {
	_, err := envrewrite.New([]envrewrite.Rule{{From: "", To: "NEW_KEY"}})
	if err == nil {
		t.Fatal("expected error for empty From, got nil")
	}
}

func TestApply_RenamesKey(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{{From: "OLD_KEY", To: "NEW_KEY"}})
	out := rw.Apply(map[string]string{"OLD_KEY": "value", "OTHER": "x"})
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if out["OTHER"] != "x" {
		t.Error("unmatched key should pass through")
	}
}

func TestApply_DropsKey(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{{From: "DROP_ME", To: ""}})
	out := rw.Apply(map[string]string{"DROP_ME": "secret", "KEEP": "yes"})
	if _, ok := out["DROP_ME"]; ok {
		t.Error("DROP_ME should have been dropped")
	}
	if out["KEEP"] != "yes" {
		t.Error("KEEP should be present")
	}
}

func TestApply_UnmatchedPassThrough(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{{From: "A", To: "B"}})
	out := rw.Apply(map[string]string{"X": "1", "Y": "2"})
	if len(out) != 2 || out["X"] != "1" || out["Y"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{{From: "A", To: "B"}})
	out := rw.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{
		{From: "DB_PASS", To: "DATABASE_PASSWORD"},
		{From: "API_KEY", To: ""},
	})
	out := rw.Apply(map[string]string{
		"DB_PASS": "secret",
		"API_KEY": "key123",
		"HOST":    "localhost",
	})
	if out["DATABASE_PASSWORD"] != "secret" {
		t.Error("expected DATABASE_PASSWORD to be set")
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should be dropped")
	}
	if out["HOST"] != "localhost" {
		t.Error("HOST should pass through")
	}
}

func TestKeys_ReturnsFromFields(t *testing.T) {
	rw, _ := envrewrite.New([]envrewrite.Rule{
		{From: "A", To: "B"},
		{From: "C", To: "D"},
	})
	keys := rw.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}
