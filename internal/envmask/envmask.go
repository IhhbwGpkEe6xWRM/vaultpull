// Package envmask provides selective masking of environment variable values
// based on key patterns, replacing sensitive values with a configurable placeholder.
package envmask

import (
	"regexp"
	"strings"
)

// DefaultPlaceholder is used when no custom placeholder is set.
const DefaultPlaceholder = "***"

// defaultPatterns are key patterns that trigger masking by default.
var defaultPatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)apikey",
	"(?i)api_key",
	"(?i)private",
	"(?i)credential",
}

// Masker selectively masks values whose keys match sensitive patterns.
type Masker struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// Option configures a Masker.
type Option func(*Masker)

// WithPlaceholder sets a custom placeholder string.
func WithPlaceholder(p string) Option {
	return func(m *Masker) { m.placeholder = p }
}

// WithPatterns replaces the default key patterns with the provided ones.
func WithPatterns(patterns []string) Option {
	return func(m *Masker) {
		m.patterns = compilePatterns(patterns)
	}
}

// New creates a Masker with default sensitive key patterns.
func New(opts ...Option) (*Masker, error) {
	m := &Masker{
		patterns:    compilePatterns(defaultPatterns),
		placeholder: DefaultPlaceholder,
	}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, re := range m.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Apply returns a copy of src with sensitive values replaced by the placeholder.
func (m *Masker) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if m.IsSensitive(k) {
			out[k] = m.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// MaskedKeys returns a sorted list of keys that would be masked.
func (m *Masker) MaskedKeys(src map[string]string) []string {
	var keys []string
	for k := range src {
		if m.IsSensitive(k) {
			keys = append(keys, k)
		}
	}
	sortStrings(keys)
	return keys
}

func compilePatterns(patterns []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			out = append(out, re)
		}
	}
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && strings.ToLower(s[j]) < strings.ToLower(s[j-1]); j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
