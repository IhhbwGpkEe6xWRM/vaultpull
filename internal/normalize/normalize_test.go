package normalize_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/normalize"
)

func TestApply_DefaultUppercase(t *testing.T) {
	n := normalize.New()
	out := n.Apply(map[string]string{"db_host": "localhost"})
	if out["DB_HOST"] != "localhost" {
		t.Fatalf("expected DB_HOST, got %v", out)
	}
}

func TestApply_HyphensReplacedWithUnderscore(t *testing.T) {
	n := normalize.New()
	out := n.Apply(map[string]string{"my-secret-key": "val"})
	if _, ok := out["MY_SECRET_KEY"]; !ok {
		t.Fatalf("expected MY_SECRET_KEY in %v", out)
	}
}

func TestApply_WithPrefix(t *testing.T) {
	n := normalize.New(normalize.WithPrefix("APP"))
	out := n.Apply(map[string]string{"token": "abc"})
	if out["APP_TOKEN"] != "abc" {
		t.Fatalf("expected APP_TOKEN, got %v", out)
	}
}

func TestApply_PrefixTrimsUnderscores(t *testing.T) {
	n := normalize.New(normalize.WithPrefix("__APP__"))
	out := n.Apply(map[string]string{"key": "v"})
	if _, ok := out["APP_KEY"]; !ok {
		t.Fatalf("expected APP_KEY, got %v", out)
	}
}

func TestApply_NoUppercase(t *testing.T) {
	n := normalize.New(normalize.WithUppercase(false))
	out := n.Apply(map[string]string{"MyKey": "v"})
	if _, ok := out["MyKey"]; !ok {
		t.Fatalf("expected MyKey unchanged, got %v", out)
	}
}

func TestApply_SpecialCharsStripped(t *testing.T) {
	n := normalize.New()
	out := n.Apply(map[string]string{"key.name/path": "v"})
	if _, ok := out["KEY_NAME_PATH"]; !ok {
		t.Fatalf("expected KEY_NAME_PATH, got %v", out)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	n := normalize.New()
	out := n.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	n := normalize.New()
	input := map[string]string{"original": "val"}
	n.Apply(input)
	if _, ok := input["original"]; !ok {
		t.Fatal("input map was mutated")
	}
}
