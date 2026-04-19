package envdiff_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envdiff"
)

func writeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCompare_AllAdded_WhenFileMissing(t *testing.T) {
	res, err := envdiff.Compare("/nonexistent/.env", map[string]string{"A": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 || res.Added[0] != "A" {
		t.Fatalf("expected A added, got %v", res.Added)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	p := writeEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := envdiff.Compare(p, map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(res.Unchanged) != 2 {
		t.Fatalf("expected 2 unchanged, got %d", len(res.Unchanged))
	}
}

func TestCompare_Modified(t *testing.T) {
	p := writeEnv(t, "KEY=old\n")
	res, err := envdiff.Compare(p, map[string]string{"KEY": "new"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Modified) != 1 || res.Modified[0] != "KEY" {
		t.Fatalf("expected KEY modified, got %v", res.Modified)
	}
}

func TestCompare_Removed(t *testing.T) {
	p := writeEnv(t, "OLD=val\n")
	res, err := envdiff.Compare(p, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "OLD" {
		t.Fatalf("expected OLD removed, got %v", res.Removed)
	}
}

func TestCompare_Mixed(t *testing.T) {
	p := writeEnv(t, "KEEP=same\nCHANGED=old\nGONE=bye\n")
	incoming := map[string]string{"KEEP": "same", "CHANGED": "new", "FRESH": "hi"}
	res, err := envdiff.Compare(p, incoming)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Added) != 1 || res.Added[0] != "FRESH" {
		t.Fatalf("added: %v", res.Added)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "GONE" {
		t.Fatalf("removed: %v", res.Removed)
	}
	if len(res.Modified) != 1 || res.Modified[0] != "CHANGED" {
		t.Fatalf("modified: %v", res.Modified)
	}
}

func TestHasChanges_FalseWhenEmpty(t *testing.T) {
	var r envdiff.Result
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
}
