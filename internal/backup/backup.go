package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Store manages local backups of .env files before overwriting them.
type Store struct {
	dir string
}

// NewStore creates a new Store that saves backups under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("backup: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save copies src into the backup directory with a timestamp suffix.
// Returns the path of the created backup file.
func (s *Store) Save(src string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // nothing to back up
		}
		return "", fmt.Errorf("backup: read source: %w", err)
	}

	base := filepath.Base(src)
	stamp := time.Now().UTC().Format("20060102T150405Z")
	dst := filepath.Join(s.dir, fmt.Sprintf("%s.%s.bak", base, stamp))

	if err := os.WriteFile(dst, data, 0600); err != nil {
		return "", fmt.Errorf("backup: write backup: %w", err)
	}
	return dst, nil
}

// List returns all backup files for the given source filename.
func (s *Store) List(src string) ([]string, error) {
	base := filepath.Base(src)
	pattern := filepath.Join(s.dir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("backup: list: %w", err)
	}
	return matches, nil
}

// Restore copies the given backup file back to dst.
func (s *Store) Restore(backupPath, dst string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("backup: read backup: %w", err)
	}
	if err := os.WriteFile(dst, data, 0600); err != nil {
		return fmt.Errorf("backup: restore write: %w", err)
	}
	return nil
}
