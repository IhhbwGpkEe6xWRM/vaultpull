// Package envdrop removes keys from a secret map based on exact names or glob patterns.
package envdrop

import (
	"fmt"
	"path"
	"sort"
)

// Dropper removes keys from a map of secrets.
type Dropper struct {
	patterns []string
}

// Option configures a Dropper.
type Option func(*Dropper)

// New creates a Dropper that will remove keys matching any of the given
// patterns. Patterns follow the same syntax as path.Match (e.g. "SECRET_*").
// Returns an error if any pattern is syntactically invalid.
func New(patterns []string) (*Dropper, error) {
	for _, p := range patterns {
		if _, err := path.Match(p, ""); err != nil {
			return nil, fmt.Errorf("envdrop: invalid pattern %q: %w", p, err)
		}
	}
	copied := make([]string, len(patterns))
	copy(copied, patterns)
	return &Dropper{patterns: copied}, nil
}

// Apply returns a new map with all keys matching any registered pattern removed.
// The original map is not mutated.
func (d *Dropper) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !d.matches(k) {
			out[k] = v
		}
	}
	return out
}

// Dropped returns the sorted list of keys from secrets that would be removed.
func (d *Dropper) Dropped(secrets map[string]string) []string {
	var keys []string
	for k := range secrets {
		if d.matches(k) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func (d *Dropper) matches(key string) bool {
	for _, p := range d.patterns {
		ok, _ := path.Match(p, key)
		if ok {
			return true
		}
	}
	return false
}
