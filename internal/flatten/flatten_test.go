package flatten_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/flatten"
)

func TestFlatten_FlatMap(t *testing.T) {
	f := flatten.New()
	input := map[string]any{"FOO": "bar", "BAZ": "qux"}
	out := f.Flatten(input)
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	f := flatten.New()
	input := map[string]any{
		"DB": map[string]any{
			"HOST": "localhost",
			"PORT": "5432",
		},
	}
	out := f.Flatten(input)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlatten_DeeplyNested(t *testing.T) {
	f := flatten.New()
	input := map[string]any{
		"A": map[string]any{
			"B": map[string]any{
				"C": "deep",
			},
		},
	}
	out := f.Flatten(input)
	if out["A_B_C"] != "deep" {
		t.Errorf("expected A_B_C=deep, got %q", out["A_B_C"])
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := flatten.NewWithSeparator(".")
	input := map[string]any{
		"db": map[string]any{"host": "localhost"},
	}
	out := f.Flatten(input)
	if out["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %q", out["db.host"])
	}
}

func TestFlatten_NonStringValue(t *testing.T) {
	f := flatten.New()
	input := map[string]any{"PORT": 8080}
	out := f.Flatten(input)
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
}

func TestFlatten_EmptyInput(t *testing.T) {
	f := flatten.New()
	out := f.Flatten(map[string]any{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestNewWithSeparator_EmptyDefaultsToUnderscore(t *testing.T) {
	f := flatten.NewWithSeparator("")
	input := map[string]any{"A": map[string]any{"B": "v"}}
	out := f.Flatten(input)
	if out["A_B"] != "v" {
		t.Errorf("expected A_B=v, got %v", out)
	}
}
