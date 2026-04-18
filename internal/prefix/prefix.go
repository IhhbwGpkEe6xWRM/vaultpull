// Package prefix provides path prefix stripping for secret keys read from Vault.
// It allows consumers to remove a common path prefix from secret keys before
// writing them to .env files, keeping output clean and portable.
package prefix

import "strings"

// Stripper removes a configured prefix from secret map keys.
type Stripper struct {
	prefix string
}

// New returns a Stripper that removes the given prefix from keys.
// Leading and trailing slashes are trimmed from the prefix.
func New(prefix string) *Stripper {
	return &Stripper{
		prefix: strings.Trim(prefix, "/"),
	}
}

// Strip returns a new map with the configured prefix removed from each key.
// Keys that do not start with the prefix are kept unchanged.
func (s *Stripper) Strip(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[s.stripKey(k)] = v
	}
	return out
}

// stripKey removes the prefix from a single key.
func (s *Stripper) stripKey(key string) string {
	if s.prefix == "" {
		return key
	}
	trimmed := strings.TrimPrefix(key, s.prefix+"/")
	if trimmed == key {
		// exact match with no trailing slash
		trimmed = strings.TrimPrefix(key, s.prefix)
	}
	return trimmed
}

// HasPrefix reports whether any key in secrets starts with the configured prefix.
func (s *Stripper) HasPrefix(secrets map[string]string) bool {
	if s.prefix == "" {
		return false
	}
	for k := range secrets {
		if strings.HasPrefix(k, s.prefix) {
			return true
		}
	}
	return false
}
