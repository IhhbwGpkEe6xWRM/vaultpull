package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func tempSnapshotFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "snapshot.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempSnapshotFile(t)
	store := snapshot.NewStore(path)

	snap := snapshot.Snapshot{
		Path:      "secret/myapp",
		Namespace: "production",
		Secrets:   map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"},
	}

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded == nil {
		t.Fatal("Load() returned nil, want snapshot")
	}
	if loaded.Path != snap.Path {
		t.Errorf("Path = %q, want %q", loaded.Path, snap.Path)
	}
	if loaded.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", loaded.Secrets["DB_HOST"], "localhost")
	}
}

func TestLoad_MissingFile_ReturnsNil(t *testing.T) {
	store := snapshot.NewStore("/nonexistent/path/snapshot.json")
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if loaded != nil {
		t.Errorf("Load() = %v, want nil", loaded)
	}
}

func TestExists_ReturnsFalseWhenMissing(t *testing.T) {
	store := snapshot.NewStore("/nonexistent/snapshot.json")
	if store.Exists() {
		t.Error("Exists() = true, want false")
	}
}

func TestExists_ReturnsTrueAfterSave(t *testing.T) {
	path := tempSnapshotFile(t)
	store := snapshot.NewStore(path)

	if err := store.Save(snapshot.Snapshot{Secrets: map[string]string{}}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if !store.Exists() {
		t.Error("Exists() = false, want true")
	}
}

func TestSave_SetsFilePermissions(t *testing.T) {
	path := tempSnapshotFile(t)
	store := snapshot.NewStore(path)

	if err := store.Save(snapshot.Snapshot{Secrets: map[string]string{"KEY": "val"}}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file perm = %o, want 0600", perm)
	}
}
