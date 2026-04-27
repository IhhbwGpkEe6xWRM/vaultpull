// Package envtrace records the origin of each secret key as it flows
// through the vaultpull pipeline. Each entry captures the Vault path the
// key was read from, making it easy to audit where a value came from.
package envtrace

import (
	"fmt"
	"sort"
	"strings"
)

// Entry describes the provenance of a single secret key.
type Entry struct {
	Key  string
	Path string
}

// Tracer records origin information for secret keys.
type Tracer struct {
	entries map[string]Entry
}

// New returns an empty Tracer.
func New() *Tracer {
	return &Tracer{entries: make(map[string]Entry)}
}

// Record associates key with the Vault path it was read from.
// If key is empty, Record returns an error.
func (t *Tracer) Record(key, path string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("envtrace: key must not be empty")
	}
	t.entries[key] = Entry{Key: key, Path: path}
	return nil
}

// RecordAll calls Record for every key in secrets using the given path.
// It returns the first error encountered, if any.
func (t *Tracer) RecordAll(secrets map[string]string, path string) error {
	for k := range secrets {
		if err := t.Record(k, path); err != nil {
			return err
		}
	}
	return nil
}

// Lookup returns the Entry for key and whether it was found.
func (t *Tracer) Lookup(key string) (Entry, bool) {
	e, ok := t.entries[key]
	return e, ok
}

// Entries returns all recorded entries sorted by key.
func (t *Tracer) Entries() []Entry {
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})
	return out
}

// Len returns the number of recorded entries.
func (t *Tracer) Len() int { return len(t.entries) }
