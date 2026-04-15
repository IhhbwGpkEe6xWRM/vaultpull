// Package filter implements namespace-based filtering for HashiCorp Vault secret paths.
//
// It provides a Matcher type that can be configured with a namespace prefix and
// used to determine whether a given Vault path belongs to that namespace, strip
// the namespace prefix from paths, or filter a slice of paths down to only those
// that fall within the namespace.
//
// Example usage:
//
//	m := filter.NewMatcher("team/backend")
//
//	// Check if a path belongs to the namespace
//	if m.Match("team/backend/database") {
//		// process secret
//	}
//
//	// Strip namespace prefix for local key naming
//	relative := m.StripNamespace("team/backend/database") // → "database"
//
//	// Filter a list of paths
//	matched := m.FilterPaths(allPaths)
package filter
