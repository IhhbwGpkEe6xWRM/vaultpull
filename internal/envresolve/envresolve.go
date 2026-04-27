// Package envresolve resolves environment variable references across multiple
// sources, applying a priority-ordered lookup chain to substitute placeholders
// in the form ${KEY} or $KEY with their resolved values.
package envresolve

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Resolver resolves variable references in a map of env values using one or
// more source maps consulted in order (first match wins).
type Resolver struct {
	sources []map[string]string
	missing string
}

// Option configures a Resolver.
type Option func(*Resolver)

// WithMissingPlaceholder sets the string substituted when a referenced key
// cannot be found in any source. Defaults to an empty string.
func WithMissingPlaceholder(placeholder string) Option {
	return func(r *Resolver) { r.missing = placeholder }
}

// New creates a Resolver that looks up references in the provided sources.
// Sources are consulted in order; the first non-empty match wins.
func New(sources []map[string]string, opts ...Option) (*Resolver, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("envresolve: at least one source map is required")
	}
	r := &Resolver{sources: sources}
	for _, o := range opts {
		o(r)
	}
	return r, nil
}

// Resolve substitutes all variable references in the values of m and returns
// a new map with the resolved values. The input map is never mutated.
func (r *Resolver) Resolve(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.expand(v)
	}
	return out
}

// ContainsReferences reports whether s contains any resolvable references.
func ContainsReferences(s string) bool {
	return refPattern.MatchString(s)
}

func (r *Resolver) expand(s string) string {
	return refPattern.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimPrefix(strings.TrimPrefix(strings.Trim(match, "${}"), "${"), "$")
		key = strings.TrimSuffix(key, "}")
		for _, src := range r.sources {
			if val, ok := src[key]; ok && val != "" {
				return val
			}
		}
		return r.missing
	})
}
