// Package pin provides secret version pinning for vaultpull.
// It allows specific Vault secret paths to be locked to a known
// version, preventing unintended updates during sync.
package pin

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a pinned version for a secret path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by,omitempty"`
}

// Store persists pin entries to a JSON file.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Entry
}

// NewStore loads or creates a pin store at the given file path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, data: make(map[string]Entry)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

// Set pins a path to a specific version.
func (s *Store) Set(path string, version int, pinnedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[path] = Entry{Path: path, Version: version, PinnedAt: time.Now().UTC(), PinnedBy: pinnedBy}
	return s.save()
}

// Get returns the pin entry for a path, and whether it exists.
func (s *Store) Get(path string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[path]
	return e, ok
}

// Remove deletes a pin for the given path.
func (s *Store) Remove(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, path)
	return s.save()
}

func (s *Store) save() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.data)
}
