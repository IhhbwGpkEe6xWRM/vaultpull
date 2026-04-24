// Package envindex builds a reverse-lookup index from secret values back to
// their keys, enabling fast duplicate-value detection across a secrets map.
package envindex

import "sort"

// Index maps a normalised value fingerprint to the list of keys that share it.
type Index struct {
	byValue map[string][]string
}

// New builds an Index from the provided secrets map.
// Keys whose values are empty strings are skipped.
func New(secrets map[string]string) *Index {
	idx := &Index{byValue: make(map[string][]string, len(secrets))}
	for k, v := range secrets {
		if v == "" {
			continue
		}
		idx.byValue[v] = append(idx.byValue[v], k)
	}
	// Sort key slices so results are deterministic.
	for v := range idx.byValue {
		sort.Strings(idx.byValue[v])
	}
	return idx
}

// KeysForValue returns all keys whose value equals v.
// Returns nil when no keys share that value.
func (idx *Index) KeysForValue(v string) []string {
	return idx.byValue[v]
}

// Duplicates returns a map of value → sorted key list for every value that
// appears under more than one key.
func (idx *Index) Duplicates() map[string][]string {
	out := make(map[string][]string)
	for v, keys := range idx.byValue {
		if len(keys) > 1 {
			out[v] = keys
		}
	}
	return out
}

// HasDuplicates reports whether any value is shared by more than one key.
func (idx *Index) HasDuplicates() bool {
	for _, keys := range idx.byValue {
		if len(keys) > 1 {
			return true
		}
	}
	return false
}

// Len returns the number of distinct non-empty values in the index.
func (idx *Index) Len() int {
	return len(idx.byValue)
}
