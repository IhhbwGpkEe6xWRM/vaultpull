// Package envsort provides deterministic ordering for secret maps.
// Keys can be sorted alphabetically or by a custom priority list,
// ensuring stable output across repeated syncs.
package envsort

import "sort"

// Sorter orders secret map keys for deterministic output.
type Sorter struct {
	priority []string
}

// Option configures a Sorter.
type Option func(*Sorter)

// WithPriority sets keys that should appear first, in order.
func WithPriority(keys ...string) Option {
	return func(s *Sorter) {
		s.priority = keys
	}
}

// New returns a Sorter with the given options.
func New(opts ...Option) *Sorter {
	s := &Sorter{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Keys returns the keys of m in sorted order, priority keys first.
func (s *Sorter) Keys(m map[string]string) []string {
	priSet := make(map[string]int, len(s.priority))
	for i, k := range s.priority {
		priSet[k] = i
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		pi, iok := priSet[keys[i]]
		pj, jok := priSet[keys[j]]
		switch {
		case iok && jok:
			return pi < pj
		case iok:
			return true
		case jok:
			return false
		default:
			return keys[i] < keys[j]
		}
	})
	return keys
}

// Apply returns a new slice of key=value pairs in sorted order.
func (s *Sorter) Apply(m map[string]string) []string {
	keys := s.Keys(m)
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = k + "=" + m[k]
	}
	return out
}
