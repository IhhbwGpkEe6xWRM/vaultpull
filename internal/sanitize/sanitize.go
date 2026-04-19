// Package sanitize provides utilities for cleaning and normalising secret
// values before they are written to .env files.
package sanitize

import (
	"strings"
	"unicode"
)

// Sanitizer cleans secret values according to a configurable policy.
type Sanitizer struct {
	stripControl bool
	trimSpace    bool
	maxLen       int
}

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithStripControl removes non-printable / control characters from values.
func WithStripControl() Option { return func(s *Sanitizer) { s.stripControl = true } }

// WithTrimSpace trims leading and trailing whitespace from values.
func WithTrimSpace() Option { return func(s *Sanitizer) { s.trimSpace = true } }

// WithMaxLen truncates values that exceed n bytes (0 = unlimited).
func WithMaxLen(n int) Option { return func(s *Sanitizer) { s.maxLen = n } }

// New returns a Sanitizer configured with the supplied options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Value applies all configured transformations to a single string.
func (s *Sanitizer) Value(v string) string {
	if s.stripControl {
		v = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) && r != '\t' {
				return -1
			}
			return r
		}, v)
	}
	if s.trimSpace {
		v = strings.TrimSpace(v)
	}
	if s.maxLen > 0 && len(v) > s.maxLen {
		v = v[:s.maxLen]
	}
	return v
}

// Map applies Value to every entry in m and returns a new map.
func (s *Sanitizer) Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = s.Value(v)
	}
	return out
}
