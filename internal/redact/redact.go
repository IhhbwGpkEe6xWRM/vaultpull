// Package redact provides utilities for scrubbing sensitive keys
// from secret maps before they are written to disk or logged.
package redact

import "strings"

// DefaultSensitivePatterns is the list of key substrings treated as sensitive.
var DefaultSensitivePatterns = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"private_key", "privatekey", "credential", "auth",
}

// Redactor scrubs values whose keys match sensitive patterns.
type Redactor struct {
	patterns []string
	placeholder string
}

// New returns a Redactor using DefaultSensitivePatterns.
func New() *Redactor {
	return NewWithPatterns(DefaultSensitivePatterns, "[REDACTED]")
}

// NewWithPatterns returns a Redactor with custom patterns and placeholder.
func NewWithPatterns(patterns []string, placeholder string) *Redactor {
	norm := make([]string, len(patterns))
	for i, p := range patterns {
		norm[i] = strings.ToLower(p)
	}
	return &Redactor{patterns: norm, placeholder: placeholder}
}

// IsSensitive reports whether key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of m with sensitive values replaced by the placeholder.
func (r *Redactor) Redact(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}
