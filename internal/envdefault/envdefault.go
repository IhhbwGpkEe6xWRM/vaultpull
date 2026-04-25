// Package envdefault provides a mechanism for applying default values to
// environment variable maps. Keys that are missing or have empty values
// in the target map are filled in from the defaults set.
package envdefault

import "fmt"

// DefaultEntry holds a key and its default value.
type DefaultEntry struct {
	Key   string
	Value string
}

// Applier applies a set of default values to an env map.
type Applier struct {
	defaults []DefaultEntry
	overwrite bool
}

// Option configures an Applier.
type Option func(*Applier)

// WithOverwrite causes Apply to overwrite existing non-empty values.
func WithOverwrite() Option {
	return func(a *Applier) {
		a.overwrite = true
	}
}

// New creates an Applier from a slice of "KEY=value" pairs.
// Malformed pairs (missing "=") are rejected with an error.
func New(pairs []string, opts ...Option) (*Applier, error) {
	a := &Applier{}
	for _, opt := range opts {
		opt(a)
	}
	for _, p := range pairs {
		for i, c := range p {
			if c == '=' {
				key := p[:i]
				val := p[i+1:]
				if key == "" {
					return nil, fmt.Errorf("envdefault: empty key in pair %q", p)
				}
				a.defaults = append(a.defaults, DefaultEntry{Key: key, Value: val})
				goto next
			}
		}
		return nil, fmt.Errorf("envdefault: malformed pair %q: missing '='", p)
	next:
	}
	return a, nil
}

// Apply returns a new map based on src with defaults filled in.
// If overwrite is false (default), only missing or empty keys are set.
// The original map is never mutated.
func (a *Applier) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	for _, d := range a.defaults {
		existing, ok := out[d.Key]
		if !ok || existing == "" || a.overwrite {
			out[d.Key] = d.Value
		}
	}
	return out
}

// Keys returns the list of keys that have registered defaults.
func (a *Applier) Keys() []string {
	keys := make([]string, len(a.defaults))
	for i, d := range a.defaults {
		keys[i] = d.Key
	}
	return keys
}
