// Package envversion tracks version numbers for secret paths, allowing
// callers to detect when a secret has been updated since last sync.
package envversion

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry holds version metadata for a single secret path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Store persists version entries to a JSON file on disk.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// NewStore loads an existing version file or returns an empty store.
func NewStore(path string) (*Store, error) {
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Set records the version for a secret path and persists to disk.
func (s *Store) Set(path string, version int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[path] = Entry{
		Path:      path,
		Version:   version,
		FetchedAt: time.Now().UTC(),
	}
	return s.save()
}

// Get returns the stored entry for a path, or false if not present.
func (s *Store) Get(path string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[path]
	return e, ok
}

// IsNewer reports whether version is strictly greater than the stored version.
// Returns true when no entry exists (treat unknown as needing sync).
func (s *Store) IsNewer(path string, version int) bool {
	e, ok := s.Get(path)
	if !ok {
		return true
	}
	return version > e.Version
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
