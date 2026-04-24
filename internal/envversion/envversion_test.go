package envversion_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envversion"
)

func tempVersionFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "versions.json")
}

func TestNewStore_EmptyWhenMissing(t *testing.T) {
	s, err := envversion.NewStore(tempVersionFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("secret/app")
	if ok {
		t.Fatal("expected no entry for unknown path")
	}
}

func TestSet_PersistsEntry(t *testing.T) {
	path := tempVersionFile(t)
	s, _ := envversion.NewStore(path)
	if err := s.Set("secret/app", 3); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	// reload from disk
	s2, err := envversion.NewStore(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	e, ok := s2.Get("secret/app")
	if !ok {
		t.Fatal("expected entry after reload")
	}
	if e.Version != 3 {
		t.Errorf("version = %d, want 3", e.Version)
	}
}

func TestIsNewer_UnknownPath_ReturnsTrue(t *testing.T) {
	s, _ := envversion.NewStore(tempVersionFile(t))
	if !s.IsNewer("secret/unknown", 1) {
		t.Fatal("expected IsNewer=true for unknown path")
	}
}

func TestIsNewer_OlderVersion_ReturnsFalse(t *testing.T) {
	s, _ := envversion.NewStore(tempVersionFile(t))
	_ = s.Set("secret/app", 5)
	if s.IsNewer("secret/app", 4) {
		t.Fatal("expected IsNewer=false for older version")
	}
}

func TestIsNewer_SameVersion_ReturnsFalse(t *testing.T) {
	s, _ := envversion.NewStore(tempVersionFile(t))
	_ = s.Set("secret/app", 5)
	if s.IsNewer("secret/app", 5) {
		t.Fatal("expected IsNewer=false for same version")
	}
}

func TestIsNewer_GreaterVersion_ReturnsTrue(t *testing.T) {
	s, _ := envversion.NewStore(tempVersionFile(t))
	_ = s.Set("secret/app", 5)
	if !s.IsNewer("secret/app", 6) {
		t.Fatal("expected IsNewer=true for greater version")
	}
}

func TestNewStore_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempVersionFile(t)
	_ = os.WriteFile(path, []byte("not-json"), 0600)
	_, err := envversion.NewStore(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGet_FetchedAtPopulated(t *testing.T) {
	s, _ := envversion.NewStore(tempVersionFile(t))
	_ = s.Set("secret/db", 2)
	e, ok := s.Get("secret/db")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.FetchedAt.IsZero() {
		t.Error("expected FetchedAt to be set")
	}
}
