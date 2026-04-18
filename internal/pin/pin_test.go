package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/pin"
)

func tempPinFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestNewStore_EmptyWhenMissing(t *testing.T) {
	s, err := pin.NewStore(tempPinFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("secret/foo")
	if ok {
		t.Error("expected no entry for missing path")
	}
}

func TestSet_PersistsEntry(t *testing.T) {
	path := tempPinFile(t)
	s, _ := pin.NewStore(path)
	if err := s.Set("secret/foo", 3, "alice"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, ok := s.Get("secret/foo")
	if !ok {
		t.Fatal("expected entry after Set")
	}
	if e.Version != 3 {
		t.Errorf("version: got %d, want 3", e.Version)
	}
	if e.PinnedBy != "alice" {
		t.Errorf("pinned_by: got %q, want alice", e.PinnedBy)
	}
}

func TestSet_ReloadsFromDisk(t *testing.T) {
	path := tempPinFile(t)
	s1, _ := pin.NewStore(path)
	s1.Set("secret/bar", 7, "bob")

	s2, err := pin.NewStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("secret/bar")
	if !ok {
		t.Fatal("expected entry in reloaded store")
	}
	if e.Version != 7 {
		t.Errorf("version: got %d, want 7", e.Version)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	path := tempPinFile(t)
	s, _ := pin.NewStore(path)
	s.Set("secret/baz", 1, "")
	if err := s.Remove("secret/baz"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, ok := s.Get("secret/baz")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestNewStore_InvalidPath(t *testing.T) {
	_, err := pin.NewStore("/no/such/dir/pins.json")
	if err == nil {
		t.Error("expected error for unwritable path")
	}
}

func TestSet_TimestampSet(t *testing.T) {
	s, _ := pin.NewStore(tempPinFile(t))
	s.Set("secret/ts", 2, "")
	e, _ := s.Get("secret/ts")
	if e.PinnedAt.IsZero() {
		t.Error("expected PinnedAt to be set")
	}
}

func init() {
	// ensure tempPinFile helper uses os.MkdirTemp indirectly via t.TempDir
	_ = os.Getenv
}
