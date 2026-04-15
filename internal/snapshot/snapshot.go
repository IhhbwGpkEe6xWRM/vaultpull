// Package snapshot provides functionality for capturing and comparing
// the state of secrets at a point in time, enabling change detection
// between successive vaultpull runs.
package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of secrets for a given path.
type Snapshot struct {
	Path      string            `json:"path"`
	Namespace string            `json:"namespace"`
	Secrets   map[string]string `json:"secrets"`
	CapturedAt time.Time        `json:"captured_at"`
}

// Store manages reading and writing snapshots to disk.
type Store struct {
	filePath string
}

// NewStore creates a new Store backed by the given file path.
func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Save writes the snapshot to disk, overwriting any existing snapshot.
func (s *Store) Save(snap Snapshot) error {
	snap.CapturedAt = time.Now().UTC()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}

// Load reads the most recent snapshot from disk.
// Returns nil, nil if no snapshot exists yet.
func (s *Store) Load() (*Snapshot, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// Exists reports whether a snapshot file is present on disk.
func (s *Store) Exists() bool {
	_, err := os.Stat(s.filePath)
	return err == nil
}
