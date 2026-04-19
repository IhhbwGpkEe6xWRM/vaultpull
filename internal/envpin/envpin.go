// Package envpin allows specific env keys to be pinned to fixed values,
// preventing vault syncs from overwriting them.
package envpin

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

// Store holds pinned key-value overrides.
type Store struct {
	mu   sync.RWMutex
	path string
	pins map[string]string
}

// NewStore loads (or creates) a pin store at path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, pins: make(map[string]string)}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.pins); err != nil {
		return nil, err
	}
	return s, nil
}

// Set pins key to value and persists.
func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pins[key] = value
	return s.save()
}

// Remove unpins a key and persists.
func (s *Store) Remove(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pins, key)
	return s.save()
}

// Apply overwrites any pinned keys in secrets with their pinned values.
func (s *Store) Apply(secrets map[string]string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for k, v := range s.pins {
		out[k] = v
	}
	return out
}

// Keys returns sorted pinned keys.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.pins))
	for k := range s.pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
