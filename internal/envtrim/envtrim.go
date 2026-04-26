// Package envtrim removes keys from a secret map whose values match
// one or more trim rules (empty, whitespace-only, or a custom predicate).
package envtrim

import "strings"

// Option configures the Trimmer.
type Option func(*Trimmer)

// WithTrimEmpty removes keys whose value is the empty string.
func WithTrimEmpty() Option {
	return func(t *Trimmer) { t.trimEmpty = true }
}

// WithTrimWhitespace removes keys whose value contains only whitespace.
func WithTrimWhitespace() Option {
	return func(t *Trimmer) { t.trimWhitespace = true }
}

// WithTrimFunc removes keys for which fn(value) returns true.
func WithTrimFunc(fn func(string) bool) Option {
	return func(t *Trimmer) { t.predicates = append(t.predicates, fn) }
}

// Trimmer filters out unwanted entries from a secret map.
type Trimmer struct {
	trimEmpty      bool
	trimWhitespace bool
	predicates     []func(string) bool
}

// New constructs a Trimmer with the given options.
func New(opts ...Option) *Trimmer {
	t := &Trimmer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply returns a copy of src with matching entries removed.
// The original map is never mutated.
func (t *Trimmer) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if t.shouldTrim(v) {
			continue
		}
		out[k] = v
	}
	return out
}

// Removed returns the keys that would be removed from src.
func (t *Trimmer) Removed(src map[string]string) []string {
	var keys []string
	for k, v := range src {
		if t.shouldTrim(v) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (t *Trimmer) shouldTrim(v string) bool {
	if t.trimEmpty && v == "" {
		return true
	}
	if t.trimWhitespace && strings.TrimSpace(v) == "" {
		return true
	}
	for _, fn := range t.predicates {
		if fn(v) {
			return true
		}
	}
	return false
}
