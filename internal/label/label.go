// Package label provides tagging and filtering of secrets by user-defined labels.
package label

import "strings"

// Set holds a collection of key=value labels.
type Set map[string]string

// Tagger attaches labels to secret maps.
type Tagger struct {
	labels Set
}

// New creates a Tagger from a slice of "key=value" strings.
// Entries that do not contain "=" are ignored.
func New(pairs []string) *Tagger {
	ls := make(Set, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		ls[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return &Tagger{labels: ls}
}

// Apply merges the tagger's labels into meta, returning a new map.
// Existing keys in meta are not overwritten.
func (t *Tagger) Apply(meta map[string]string) map[string]string {
	out := make(map[string]string, len(meta)+len(t.labels))
	for k, v := range t.labels {
		out[k] = v
	}
	for k, v := range meta {
		out[k] = v
	}
	return out
}

// Matches reports whether all pairs in filter are present and equal in meta.
func Matches(meta map[string]string, filter Set) bool {
	for k, v := range filter {
		if meta[k] != v {
			return false
		}
	}
	return true
}
