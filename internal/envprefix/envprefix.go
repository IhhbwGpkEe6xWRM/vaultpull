// Package envprefix provides utilities for adding, removing, and matching
// environment variable key prefixes during secret synchronisation.
package envprefix

import "strings"

// Transformer applies a fixed prefix to every key in a secrets map and can
// strip that same prefix when reading keys back.
type Transformer struct {
	prefix string
}

// New returns a Transformer for the given prefix. Leading and trailing
// underscores are trimmed so callers do not need to worry about double
// underscores when the prefix is joined to a key.
func New(prefix string) *Transformer {
	return &Transformer{
		prefix: strings.Trim(strings.ToUpper(prefix), "_"),
	}
}

// Apply returns a new map where every key is prefixed with t.prefix + "_".
// If the prefix is empty the original map is returned unchanged.
func (t *Transformer) Apply(secrets map[string]string) map[string]string {
	if t.prefix == "" {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[t.prefix+"_"+k] = v
	}
	return out
}

// Strip returns a new map where t.prefix + "_" is removed from every key that
// carries it. Keys that do not start with the prefix are kept as-is.
func (t *Transformer) Strip(secrets map[string]string) map[string]string {
	if t.prefix == "" {
		return secrets
	}
	pfx := t.prefix + "_"
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, pfx) {
			out[strings.TrimPrefix(k, pfx)] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// HasPrefix reports whether key starts with the transformer's prefix.
func (t *Transformer) HasPrefix(key string) bool {
	if t.prefix == "" {
		return true
	}
	return strings.HasPrefix(key, t.prefix+"_")
}
