// Package multipath provides support for reading secrets from multiple
// Vault paths and merging them into a single map.
package multipath

import (
	"fmt"
	"strings"
)

// Reader is the interface for reading secrets from a single path.
type Reader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Merger reads from multiple Vault paths and merges the results.
// Later paths take precedence over earlier ones on key conflicts.
type Merger struct {
	reader Reader
	paths  []string
}

// New creates a Merger for the given ordered list of paths.
func New(reader Reader, paths []string) (*Merger, error) {
	if reader == nil {
		return nil, fmt.Errorf("multipath: reader must not be nil")
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("multipath: at least one path is required")
	}
	cleaned := make([]string, len(paths))
	for i, p := range paths {
		cleaned[i] = strings.Trim(p, "/")
	}
	return &Merger{reader: reader, paths: cleaned}, nil
}

// Merge reads all configured paths and returns a merged secret map.
// Keys from later paths overwrite keys from earlier paths.
func (m *Merger) Merge() (map[string]string, error) {
	result := make(map[string]string)
	for _, p := range m.paths {
		secrets, err := m.reader.ReadSecrets(p)
		if err != nil {
			return nil, fmt.Errorf("multipath: reading %q: %w", p, err)
		}
		for k, v := range secrets {
			result[k] = v
		}
	}
	return result, nil
}

// Paths returns the cleaned list of paths configured on the Merger.
func (m *Merger) Paths() []string {
	out := make([]string, len(m.paths))
	copy(out, m.paths)
	return out
}
