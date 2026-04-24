package envscope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envscope"
)

func newGuard(t *testing.T, prefixes []string) *envscope.Guard {
	t.Helper()
	s := newScope(t, prefixes)
	return envscope.NewGuard(s)
}

func TestWrite_AllowedKey_Succeeds(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	if err := g.Write(dst, "APP_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["APP_HOST"] != "localhost" {
		t.Errorf("expected value to be written")
	}
}

func TestWrite_BlockedKey_ReturnsError(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	err := g.Write(dst, "DB_PASS", "secret")
	if err == nil {
		t.Fatal("expected error for blocked key")
	}
	if _, ok := dst["DB_PASS"]; ok {
		t.Error("blocked key must not be written to dst")
	}
}

func TestBlocked_ReturnsSortedKeys(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	_ = g.Write(dst, "Z_KEY", "v")
	_ = g.Write(dst, "A_KEY", "v")
	blocked := g.Blocked()
	if len(blocked) != 2 {
		t.Fatalf("expected 2 blocked keys, got %d", len(blocked))
	}
	if blocked[0] != "A_KEY" || blocked[1] != "Z_KEY" {
		t.Errorf("expected sorted order, got %v", blocked)
	}
}

func TestWriteAll_MixedKeys_ReturnsError(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	src := map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "secret",
	}
	err := g.WriteAll(dst, src)
	if err == nil {
		t.Fatal("expected combined error")
	}
	if _, ok := dst["APP_HOST"]; !ok {
		t.Error("allowed key should still be written")
	}
	if _, ok := dst["DB_PASS"]; ok {
		t.Error("blocked key must not appear in dst")
	}
}

func TestWriteAll_AllAllowed_NoError(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	err := g.WriteAll(dst, map[string]string{"APP_A": "1", "APP_B": "2"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(dst) != 2 {
		t.Errorf("expected 2 keys in dst, got %d", len(dst))
	}
}

func TestReset_ClearsBlockedHistory(t *testing.T) {
	g := newGuard(t, []string{"APP"})
	dst := map[string]string{}
	_ = g.Write(dst, "DB_PASS", "v")
	g.Reset()
	if len(g.Blocked()) != 0 {
		t.Error("expected blocked list to be empty after Reset")
	}
}
