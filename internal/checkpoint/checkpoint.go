// Package checkpoint tracks the last successful sync time and secret path
// so that incremental syncs can skip unchanged paths.
package checkpoint

import (
	"encoding/json"
	"os"
	"time"
)

// Entry holds the recorded state of a single sync operation.
type Entry struct {
	Path      string    `json:"path"`
	SyncedAt  time.Time `json:"synced_at"`
	Checksum  string    `json:"checksum"`
}

// Store persists checkpoint entries to a JSON file.
type Store struct {
	path    string
	entries map[string]Entry
}

// NewStore opens or creates a checkpoint file at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Set records a checkpoint for the given secret path.
func (s *Store) Set(secretPath, checksum string, at time.Time) error {
	s.entries[secretPath] = Entry{Path: secretPath, SyncedAt: at, Checksum: checksum}
	return s.flush()
}

// Get returns the checkpoint entry for the given secret path, if any.
func (s *Store) Get(secretPath string) (Entry, bool) {
	e, ok := s.entries[secretPath]
	return e, ok
}

// IsFresh returns true when the stored checksum matches and the entry is
// younger than maxAge.
func (s *Store) IsFresh(secretPath, checksum string, maxAge time.Duration, now time.Time) bool {
	e, ok := s.entries[secretPath]
	if !ok {
		return false
	}
	return e.Checksum == checksum && now.Sub(e.SyncedAt) < maxAge
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
