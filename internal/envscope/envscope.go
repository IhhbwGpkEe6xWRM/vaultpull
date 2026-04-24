// Package envscope restricts which environment variable keys are visible
// or writable based on a declared scope (allow-list of key prefixes).
package envscope

import (
	"fmt"
	"strings"
)

// Scope enforces a prefix-based allow-list over a map of env vars.
type Scope struct {
	prefixes []string
}

// New creates a Scope from a list of key prefixes. Prefixes are uppercased
// and trailing underscores are trimmed for normalisation.
func New(prefixes []string) (*Scope, error) {
	var cleaned []string
	for _, p := range prefixes {
		p = strings.TrimRight(strings.ToUpper(strings.TrimSpace(p)), "_")
		if p == "" {
			continue
		}
		cleaned = append(cleaned, p)
	}
	return &Scope{prefixes: cleaned}, nil
}

// Allows reports whether key is permitted by the scope.
// If no prefixes are configured every key is allowed.
func (s *Scope) Allows(key string) bool {
	if len(s.prefixes) == 0 {
		return true
	}
	upper := strings.ToUpper(key)
	for _, p := range s.prefixes {
		if upper == p || strings.HasPrefix(upper, p+"_") {
			return true
		}
	}
	return false
}

// Filter returns a copy of m containing only keys permitted by the scope.
func (s *Scope) Filter(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if s.Allows(k) {
			out[k] = v
		}
	}
	return out
}

// Validate returns an error listing any keys in m that fall outside the scope.
func (s *Scope) Validate(m map[string]string) error {
	var violations []string
	for k := range m {
		if !s.Allows(k) {
			violations = append(violations, k)
		}
	}
	if len(violations) == 0 {
		return nil
	}
	return fmt.Errorf("envscope: keys outside scope: %s", strings.Join(violations, ", "))
}
