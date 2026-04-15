// Package filter provides namespace-based filtering utilities for Vault secret paths.
package filter

import (
	"strings"
)

// Matcher determines whether a given secret path matches a namespace filter.
type Matcher struct {
	namespace string
}

// NewMatcher creates a new Matcher for the given namespace.
// An empty namespace matches all paths.
func NewMatcher(namespace string) *Matcher {
	ns := strings.Trim(namespace, "/")
	return &Matcher{namespace: ns}
}

// Match returns true if the path belongs to the configured namespace.
// If the namespace is empty, all paths match.
func (m *Matcher) Match(path string) bool {
	if m.namespace == "" {
		return true
	}
	clean := strings.Trim(path, "/")
	return clean == m.namespace ||
		strings.HasPrefix(clean, m.namespace+"/")
}

// StripNamespace removes the namespace prefix from a path, returning
// only the relative portion. If the path does not match, it is returned unchanged.
func (m *Matcher) StripNamespace(path string) string {
	if m.namespace == "" {
		return strings.Trim(path, "/")
	}
	clean := strings.Trim(path, "/")
	prefix := m.namespace + "/"
	if strings.HasPrefix(clean, prefix) {
		return clean[len(prefix):]
	}
	return clean
}

// FilterPaths returns only the paths from the provided slice that match the namespace.
func (m *Matcher) FilterPaths(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		if m.Match(p) {
			result = append(result, p)
		}
	}
	return result
}
