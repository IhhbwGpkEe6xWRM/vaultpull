package envlock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envlock"
)

func tempEnvFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(time.Second)

	lk, err := l.Acquire(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer lk.Release()

	lockPath := filepath.Join(filepath.Dir(path), ".env.lock")
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("expected lock file to exist")
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(time.Second)

	lk, err := l.Acquire(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := lk.Release(); err != nil {
		t.Fatalf("release error: %v", err)
	}

	lockPath := filepath.Join(filepath.Dir(path), ".env.lock")
	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Error("expected lock file to be removed")
	}
}

func TestAcquire_TimeoutWhenAlreadyLocked(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(100 * time.Millisecond)

	lk, err := l.Acquire(path)
	if err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer lk.Release()

	_, err = l.Acquire(path)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestAcquire_SucceedsAfterRelease(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(time.Second)

	lk, _ := l.Acquire(path)
	_ = lk.Release()

	lk2, err := l.Acquire(path)
	if err != nil {
		t.Fatalf("expected success after release, got: %v", err)
	}
	_ = lk2.Release()
}

func TestRelease_IdempotentOnMissingFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(time.Second)

	lk, _ := l.Acquire(path)
	_ = lk.Release()

	if err := lk.Release(); err != nil {
		t.Errorf("second release should not error, got: %v", err)
	}
}
