package backup_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/backup"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "backup-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSave_ReturnsEmptyWhenSourceMissing(t *testing.T) {
	store, _ := backup.NewStore(tempDir(t))
	path, err := store.Save("/nonexistent/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestSave_CreatesBackupFile(t *testing.T) {
	src := filepath.Join(tempDir(t), ".env")
	_ = os.WriteFile(src, []byte("KEY=val\n"), 0600)

	store, _ := backup.NewStore(tempDir(t))
	path, err := store.Save(src)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !strings.HasSuffix(path, ".bak") {
		t.Errorf("expected .bak suffix, got %q", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("backup file not found: %v", err)
	}
}

func TestSave_PreservesContent(t *testing.T) {
	src := filepath.Join(tempDir(t), ".env")
	content := []byte("SECRET=abc123\n")
	_ = os.WriteFile(src, content, 0600)

	store, _ := backup.NewStore(tempDir(t))
	path, _ := store.Save(src)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q want %q", got, content)
	}
}

func TestList_ReturnsBackups(t *testing.T) {
	bdir := tempDir(t)
	src := filepath.Join(tempDir(t), ".env")
	_ = os.WriteFile(src, []byte("A=1"), 0600)

	store, _ := backup.NewStore(bdir)
	store.Save(src)
	store.Save(src)

	files, err := store.List(src)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 backups, got %d", len(files))
	}
}

func TestRestore_WritesContentToDst(t *testing.T) {
	src := filepath.Join(tempDir(t), ".env")
	_ = os.WriteFile(src, []byte("ORIG=1\n"), 0600)

	store, _ := backup.NewStore(tempDir(t))
	bakPath, _ := store.Save(src)

	dst := filepath.Join(tempDir(t), ".env.restored")
	if err := store.Restore(bakPath, dst); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "ORIG=1\n" {
		t.Errorf("unexpected content: %q", data)
	}
}
