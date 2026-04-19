package envguard_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envguard"
)

func TestCheck_NoConflicts_BothUnchanged(t *testing.T) {
	g := envguard.New(map[string]string{"FOO": "bar"})
	vs, err := g.Check(
		map[string]string{"FOO": "bar"},
		map[string]string{"FOO": "bar"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheck_OnlyIncomingChanged_NoConflict(t *testing.T) {
	g := envguard.New(map[string]string{"FOO": "bar"})
	vs, _ := g.Check(
		map[string]string{"FOO": "bar"},  // local unchanged
		map[string]string{"FOO": "baz"},  // incoming changed
	)
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheck_OnlyLocalChanged_NoConflict(t *testing.T) {
	g := envguard.New(map[string]string{"FOO": "bar"})
	vs, _ := g.Check(
		map[string]string{"FOO": "local-edit"}, // local changed
		map[string]string{"FOO": "bar"},         // incoming unchanged
	)
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheck_BothChanged_Conflict(t *testing.T) {
	g := envguard.New(map[string]string{"FOO": "original"})
	vs, err := g.Check(
		map[string]string{"FOO": "local-edit"},
		map[string]string{"FOO": "vault-update"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", vs[0].Key)
	}
}

func TestCheck_MultipleConflicts_Sorted(t *testing.T) {
	snap := map[string]string{"AAA": "1", "BBB": "2", "CCC": "3"}
	g := envguard.New(snap)
	local := map[string]string{"AAA": "local-a", "BBB": "local-b", "CCC": "local-c"}
	incoming := map[string]string{"AAA": "vault-a", "BBB": "vault-b", "CCC": "vault-c"}
	vs, _ := g.Check(local, incoming)
	if len(vs) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(vs))
	}
	if vs[0].Key != "AAA" || vs[1].Key != "BBB" || vs[2].Key != "CCC" {
		t.Error("violations not sorted")
	}
}

func TestCheck_KeyNotInSnapshot_Ignored(t *testing.T) {
	g := envguard.New(map[string]string{})
	vs, _ := g.Check(
		map[string]string{"NEW": "local"},
		map[string]string{"NEW": "vault"},
	)
	if len(vs) != 0 {
		t.Fatal("new key not in snapshot should not be a conflict")
	}
}

func TestSummary_NoViolations(t *testing.T) {
	s := envguard.Summary(nil)
	if s != "no conflicts detected" {
		t.Errorf("unexpected: %s", s)
	}
}

func TestSummary_WithViolations(t *testing.T) {
	vs := []envguard.Violation{{Key: "FOO", Local: "a", Incoming: "b"}}
	s := envguard.Summary(vs)
	if !strings.Contains(s, "FOO") {
		t.Error("summary should contain key name")
	}
	if !strings.Contains(s, "1 conflict") {
		t.Error("summary should mention conflict count")
	}
}
