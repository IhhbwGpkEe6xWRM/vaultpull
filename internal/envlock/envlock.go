// Package envlock provides file-based locking to prevent concurrent writes
// to the same .env output file during a vaultpull sync operation.
package envlock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Lock represents a held file lock.
type Lock struct {
	path    string
	lockPath string
}

// Locker acquires and releases locks on env file paths.
type Locker struct {
	timeout time.Duration
	pollInterval time.Duration
}

// New returns a Locker with the given timeout for acquiring a lock.
func New(timeout time.Duration) *Locker {
	return &Locker{
		timeout:      timeout,
		pollInterval: 50 * time.Millisecond,
	}
}

// Acquire attempts to acquire a lock for the given file path.
// It returns a Lock on success or an error if the timeout is exceeded.
func (l *Locker) Acquire(path string) (*Lock, error) {
	lockPath := lockFilePath(path)
	deadline := time.Now().Add(l.timeout)

	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_ = f.Close()
			return &Lock{path: path, lockPath: lockPath}, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("envlock: open lock file: %w", err)
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("envlock: timeout acquiring lock for %s", path)
		}
		time.Sleep(l.pollInterval)
	}
}

// Release removes the lock file,eing the lock.
func (lk *Lock) Release() error {
	if err := os.Remove(lk.lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("envlock: release: %w", err)
	}
	return nil
}

func lockFilePath(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	return filepath.Join(dir, "."+base+".lock")
}
