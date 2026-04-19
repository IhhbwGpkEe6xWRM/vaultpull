package envpin

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPinFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestNewStore_EmptyWhenMissing(t *testing.T) {
	s, err := NewStore(tempPinFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Keys()) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestSet_PersistsEntry(t *testing.T) {
	path := tempPinFile(t)
	s, _ := NewStore(path)
	if err := s.Set("API_KEY", "pinned-value"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	s2, _ := NewStore(path)
	keys := s2.Keys()
	if len(keys) != 1 || keys[0] != "API_KEY" {
		t.Errorf("expected API_KEY pinned, got %v", keys)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	path := tempPinFile(t)
	s, _ := NewStore(path)
	_ = s.Set("FOO", "bar")
	if err := s.Remove("FOO"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if len(s.Keys()) != 0 {
		t.Errorf("expected empty after remove")
	}
	// verify persisted
	s2, _ := NewStore(path)
	if len(s2.Keys()) != 0 {
		t.Errorf("expected empty on reload")
	}
}

func TestApply_OverwritesPinnedKeys(t *testing.T) {
	s, _ := NewStore(tempPinFile(t))
	_ = s.Set("SECRET", "pinned")
	secrets := map[string]string{"SECRET": "from-vault", "OTHER": "value"}
	out := s.Apply(secrets)
	if out["SECRET"] != "pinned" {
		t.Errorf("expected pinned value, got %q", out["SECRET"])
	}
	if out["OTHER"] != "value" {
		t.Errorf("expected vault value for OTHER")
	}
}

func TestApply_AddsMissingPinnedKeys(t *testing.T) {
	s, _ := NewStore(tempPinFile(t))
	_ = s.Set("EXTRA", "injected")
	out := s.Apply(map[string]string{})
	if out["EXTRA"] != "injected" {
		t.Errorf("expected injected key")
	}
}

func TestNewStore_InvalidJSON(t *testing.T) {
	path := tempPinFile(t)
	_ = os.WriteFile(path, []byte("not-json"), 0600)
	_, err := NewStore(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
