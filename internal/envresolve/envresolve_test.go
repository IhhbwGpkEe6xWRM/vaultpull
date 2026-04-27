package envresolve

import (
	"testing"
)

func newResolver(t *testing.T, sources []map[string]string, opts ...Option) *Resolver {
	t.Helper()
	r, err := New(sources, opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestNew_NoSources_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestResolve_NoReferences(t *testing.T) {
	r := newResolver(t, []map[string]string{{"HOST": "localhost"}})
	out := r.Resolve(map[string]string{"KEY": "plain"})
	if out["KEY"] != "plain" {
		t.Errorf("got %q, want %q", out["KEY"], "plain")
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	src := map[string]string{"HOST": "db.internal"}
	r := newResolver(t, []map[string]string{src})
	out := r.Resolve(map[string]string{"DSN": "postgres://${HOST}/app"})
	want := "postgres://db.internal/app"
	if out["DSN"] != want {
		t.Errorf("got %q, want %q", out["DSN"], want)
	}
}

func TestResolve_NoBraceStyle(t *testing.T) {
	src := map[string]string{"PORT": "5432"}
	r := newResolver(t, []map[string]string{src})
	out := r.Resolve(map[string]string{"ADDR": "localhost:$PORT"})
	want := "localhost:5432"
	if out["ADDR"] != want {
		t.Errorf("got %q, want %q", out["ADDR"], want)
	}
}

func TestResolve_MissingKey_EmptyByDefault(t *testing.T) {
	r := newResolver(t, []map[string]string{{}})
	out := r.Resolve(map[string]string{"X": "${UNDEFINED}"})
	if out["X"] != "" {
		t.Errorf("got %q, want empty string", out["X"])
	}
}

func TestResolve_MissingKey_CustomPlaceholder(t *testing.T) {
	r := newResolver(t, []map[string]string{{}}, WithMissingPlaceholder("<unset>"))
	out := r.Resolve(map[string]string{"X": "${MISSING}"})
	if out["X"] != "<unset>" {
		t.Errorf("got %q, want %q", out["X"], "<unset>")
	}
}

func TestResolve_FirstSourceWins(t *testing.T) {
	primary := map[string]string{"TOKEN": "primary-token"}
	fallback := map[string]string{"TOKEN": "fallback-token"}
	r := newResolver(t, []map[string]string{primary, fallback})
	out := r.Resolve(map[string]string{"AUTH": "${TOKEN}"})
	if out["AUTH"] != "primary-token" {
		t.Errorf("got %q, want %q", out["AUTH"], "primary-token")
	}
}

func TestResolve_FallsBackToSecondSource(t *testing.T) {
	primary := map[string]string{}
	fallback := map[string]string{"REGION": "us-east-1"}
	r := newResolver(t, []map[string]string{primary, fallback})
	out := r.Resolve(map[string]string{"LOC": "${REGION}"})
	if out["LOC"] != "us-east-1" {
		t.Errorf("got %q, want %q", out["LOC"], "us-east-1")
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	r := newResolver(t, []map[string]string{{"A": "1"}})
	input := map[string]string{"K": "${A}"}
	r.Resolve(input)
	if input["K"] != "${A}" {
		t.Error("Resolve mutated the input map")
	}
}

func TestContainsReferences_True(t *testing.T) {
	if !ContainsReferences("hello ${WORLD}") {
		t.Error("expected true for string with reference")
	}
}

func TestContainsReferences_False(t *testing.T) {
	if ContainsReferences("plain value") {
		t.Error("expected false for string without reference")
	}
}
