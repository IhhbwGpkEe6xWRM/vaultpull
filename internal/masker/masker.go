// Package masker provides utilities for redacting sensitive secret values
// before they are written to logs, output, or audit trails.
package masker

import "strings"

const (
	// DefaultMask is the string used to replace sensitive values.
	DefaultMask = "***"

	// visibleChars is the number of characters to reveal at the start/end
	// when using partial masking.
	visibleChars = 4
)

// Masker redacts secret values for safe display.
type Masker struct {
	mask    string
	partial bool
}

// New returns a Masker that fully replaces values with the default mask.
func New() *Masker {
	return &Masker{mask: DefaultMask}
}

// NewPartial returns a Masker that reveals the first and last few characters
// of a value, useful for confirming identity without exposing the secret.
func NewPartial() *Masker {
	return &Masker{mask: DefaultMask, partial: true}
}

// Mask returns a redacted version of the given value.
func (m *Masker) Mask(value string) string {
	if value == "" {
		return value
	}
	if !m.partial || len(value) <= visibleChars*2 {
		return m.mask
	}
	return value[:visibleChars] + m.mask + value[len(value)-visibleChars:]
}

// MaskMap returns a copy of the map with all values redacted.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.Mask(v)
	}
	return out
}

// ContainsSensitive performs a naive check to see whether a string appears
// to contain a raw secret value from the provided map.
func ContainsSensitive(s string, secrets map[string]string) bool {
	for _, v := range secrets {
		if v != "" && strings.Contains(s, v) {
			return true
		}
	}
	return false
}
