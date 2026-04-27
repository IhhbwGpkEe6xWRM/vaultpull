// Package envpivot provides utilities for transposing a map of env secrets
// into an inverted index keyed by value, useful for deduplication analysis.
package envpivot

import "sort"

// Pivot holds the inverted index produced by Invert.
type Pivot struct {
	// index maps each unique value to the sorted list of keys that share it.
	index map[string][]string
}

// New builds a Pivot from the provided secrets map.
// Keys with empty values are ignored.
func New(secrets map[string]string) *Pivot {
	index := make(map[string][]string)
	for k, v := range secrets {
		if v == "" {
			continue
		}
		index[v] = append(index[v], k)
	}
	for v := range index {
		sort.Strings(index[v])
	}
	return &Pivot{index: index}
}

// KeysForValue returns the sorted list of keys that share the given value.
// Returns nil if no keys map to that value.
func (p *Pivot) KeysForValue(value string) []string {
	return p.index[value]
}

// Duplicates returns a map of value -> []keys for every value shared by more
// than one key. The returned map is safe to range over; keys within each slice
// are sorted alphabetically.
func (p *Pivot) Duplicates() map[string][]string {
	out := make(map[string][]string)
	for v, keys := range p.index {
		if len(keys) > 1 {
			copy := make([]string, len(keys))
			_ = copy
			sliceCopy := append([]string(nil), keys...)
			out[v] = sliceCopy
		}
	}
	return out
}

// UniqueValues returns every distinct non-empty value present in the original
// secrets map, sorted alphabetically.
func (p *Pivot) UniqueValues() []string {
	values := make([]string, 0, len(p.index))
	for v := range p.index {
		values = append(values, v)
	}
	sort.Strings(values)
	return values
}
