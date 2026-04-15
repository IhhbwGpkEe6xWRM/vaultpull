// Package cache provides a simple file-backed cache for Vault secrets,
// allowing vaultpull to skip re-fetching secrets that have not changed.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single cached secret snapshot.
type Entry struct {
	Path      string            `json:"path"`
	Checksum  string            `json:"checksum"`
	FetchedAt time.Time         `json:"fetched_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Store manages a JSON cache file on disk.
type Store struct {
	filePath string
	entries  map[string]Entry
}

// NewStore opens (or creates) a cache store at the given file path.
func NewStore(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
		entries:  make(map[string]Entry),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// IsFresh returns true when the cached entry for path matches the given secrets.
func (s *Store) IsFresh(path string, secrets map[string]string) bool {
	e, ok := s.entries[path]
	if !ok {
		return false
	}
	return e.Checksum == checksum(secrets)
}

// Set stores an entry for the given path.
func (s *Store) Set(path string, secrets map[string]string) {
	s.entries[path] = Entry{
		Path:      path,
		Checksum:  checksum(secrets),
		FetchedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
}

// Save flushes the in-memory entries to disk.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.entries)
}

func (s *Store) load() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.entries)
}

func checksum(secrets map[string]string) string {
	b, _ := json.Marshal(secrets)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
