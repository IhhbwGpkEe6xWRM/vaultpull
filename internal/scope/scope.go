// Package scope provides path-based scoping to restrict Vault secret reads
// to a declared set of allowed prefixes.
package scope

import "strings"

// Scope restricts secret paths to a set of allowed prefixes.
type Scope struct {
	prefixes []string
}

// New returns a Scope that allows only paths matching one of the given prefixes.
// Leading and trailing slashes are trimmed from each prefix.
func New(prefixes []string) *Scope {
	cleaned := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		p = strings.Trim(p, "/")
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return &Scope{prefixes: cleaned}
}

// Allows reports whether path is permitted by the scope.
// If no prefixes are configured, all paths are allowed.
func (s *Scope) Allows(path string) bool {
	if len(s.prefixes) == 0 {
		return true
	}
	path = strings.Trim(path, "/")
	for _, p := range s.prefixes {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// Filter returns only those paths from the input slice that are allowed.
func (s *Scope) Filter(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if s.Allows(p) {
			out = append(out, p)
		}
	}
	return out
}

// Prefixes returns the configured prefix list.
func (s *Scope) Prefixes() []string {
	return s.prefixes
}
