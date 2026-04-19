package envexpand_test

import (
	"testing"

	"github.com/hashicorp/vaultpull/internal/envexpand"
)

func TestMultiPass_ChainedReferences(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{
		"A": "${B}",
		"B": "${C}",
		"C": "resolved",
	}
	out := e.MultiPass(in, 5)
	if out["A"] != "resolved" {
		t.Fatalf("got %q, want %q", out["A"], "resolved")
	}
	if out["B"] != "resolved" {
		t.Fatalf("got %q, want %q", out["B"], "resolved")
	}
}

func TestMultiPass_StopsEarlyWhenStable(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{"FOO": "bar"}
	// Should not panic or loop even with high pass count.
	out := e.MultiPass(in, 100)
	if out["FOO"] != "bar" {
		t.Fatalf("unexpected change: %q", out["FOO"])
	}
}

func TestMultiPass_ZeroPasses_DefaultsToFive(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{
		"X": "${Y}",
		"Y": "hello",
	}
	out := e.MultiPass(in, 0)
	if out["X"] != "hello" {
		t.Fatalf("got %q, want %q", out["X"], "hello")
	}
}

func TestMultiPass_DoesNotMutateInput(t *testing.T) {
	e := envexpand.New()
	in := map[string]string{"A": "${B}", "B": "val"}
	orig := map[string]string{"A": "${A}", "B": "val"}
	_ = orig
	e.MultiPass(in, 3)
	if in["A"] != "${B}" {
		t.Fatal("input was mutated")
	}
}
