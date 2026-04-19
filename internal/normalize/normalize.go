// Package normalize provides key normalization for secret maps.
// It ensures consistent formatting of environment variable keys
// before they are written to .env files.
package normalize

import (
	"strings"
	"unicode"
)

// Normalizer transforms secret map keys into normalized env var names.
type Normalizer struct {
	prefix    string
	uppercase bool
}

// Option configures a Normalizer.
type Option func(*Normalizer)

// WithPrefix prepends a static prefix to every key.
func WithPrefix(p string) Option {
	return func(n *Normalizer) {
		n.prefix = strings.ToUpper(strings.Trim(p, "_"))
	}
}

// WithUppercase forces all keys to uppercase.
func WithUppercase(u bool) Option {
	return func(n *Normalizer) {
		n.uppercase = u
	}
}

// New returns a Normalizer configured with the given options.
// Uppercase normalization is enabled by default.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{uppercase: true}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Apply normalizes all keys in the provided map and returns a new map.
// Original map is not mutated.
func (n *Normalizer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[n.normalizeKey(k)] = v
	}
	return out
}

// normalizeKey sanitizes and transforms a single key.
func (n *Normalizer) normalizeKey(key string) string {
	var b strings.Builder
	for _, r := range key {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	result := strings.Trim(b.String(), "_")
	if n.uppercase {
		result = strings.ToUpper(result)
	}
	if n.prefix != "" {
		result = n.prefix + "_" + result
	}
	return result
}
