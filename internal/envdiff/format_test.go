package envdiff_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envdiff"
)

func TestFormat_NoChanges(t *testing.T) {
	var b strings.Builder
	envdiff.Format(&b, envdiff.Result{Unchanged: []string{"A"}}, false)
	if !strings.Contains(b.String(), "no changes") {
		t.Fatalf("unexpected output: %q", b.String())
	}
}

func TestFormat_ShowsAdded(t *testing.T) {
	var b strings.Builder
	envdiff.Format(&b, envdiff.Result{Added: []string{"NEW_KEY"}}, false)
	if !strings.Contains(b.String(), "+ NEW_KEY") {
		t.Fatalf("unexpected output: %q", b.String())
	}
}

func TestFormat_ShowsRemoved(t *testing.T) {
	var b strings.Builder
	envdiff.Format(&b, envdiff.Result{Removed: []string{"OLD_KEY"}}, false)
	if !strings.Contains(b.String(), "- OLD_KEY") {
		t.Fatalf("unexpected output: %q", b.String())
	}
}

func TestFormat_ShowsModified(t *testing.T) {
	var b strings.Builder
	envdiff.Format(&b, envdiff.Result{Modified: []string{"MOD_KEY"}}, false)
	if !strings.Contains(b.String(), "~ MOD_KEY") {
		t.Fatalf("unexpected output: %q", b.String())
	}
}

func TestFormat_SummaryLine(t *testing.T) {
	var b strings.Builder
	res := envdiff.Result{
		Added:    []string{"A"},
		Removed:  []string{"B"},
		Modified: []string{"C"},
	}
	envdiff.Format(&b, res, false)
	out := b.String()
	if !strings.Contains(out, "1 added") || !strings.Contains(out, "1 removed") || !strings.Contains(out, "1 modified") {
		t.Fatalf("unexpected summary: %q", out)
	}
}

func TestFormat_ColourAdded(t *testing.T) {
	var b strings.Builder
	envdiff.Format(&b, envdiff.Result{Added: []string{"X"}}, true)
	if !strings.Contains(b.String(), "\033[32m") {
		t.Fatal("expected green ANSI code for added key")
	}
}
