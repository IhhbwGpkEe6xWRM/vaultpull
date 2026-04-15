package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func tempCacheFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "cache.json")
}

func TestNewStore_CreatesEmptyStore(t *testing.T) {
	s, err := NewStore(tempCacheFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestIsFresh_MissingEntry(t *testing.T) {
	s, _ := NewStore(tempCacheFile(t))
	if s.IsFresh("secret/app", map[string]string{"KEY": "val"}) {
		t.Error("expected false for missing entry")
	}
}

func TestSetAndIsFresh(t *testing.T) {
	s, _ := NewStore(tempCacheFile(t))
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s.Set("secret/app", secrets)

	if !s.IsFresh("secret/app", secrets) {
		t.Error("expected entry to be fresh after Set")
	}
}

func TestIsFresh_DifferentSecrets(t *testing.T) {
	s, _ := NewStore(tempCacheFile(t))
	s.Set("secret/app", map[string]string{"FOO": "bar"})

	if s.IsFresh("secret/app", map[string]string{"FOO": "changed"}) {
		t.Error("expected false when secrets differ")
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tempCacheFile(t)
	secrets := map[string]string{"TOKEN": "abc123"}

	s1, _ := NewStore(path)
	s1.Set("secret/svc", secrets)
	if err := s1.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	s2, err := NewStore(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if !s2.IsFresh("secret/svc", secrets) {
		t.Error("expected reloaded store to recognise fresh entry")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "cache.json")
	s, _ := NewStore(path)
	s.Set("secret/x", map[string]string{"A": "1"})
	if err := s.Save(); err != nil {
		t.Fatalf("expected Save to create parent dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("cache file not found: %v", err)
	}
}

func TestChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"K": "v"}
	if checksum(secrets) != checksum(secrets) {
		t.Error("checksum should be deterministic")
	}
}
