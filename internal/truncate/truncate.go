// Package truncate provides utilities for truncating secret values
// in output to prevent accidental exposure in logs and terminals.
package truncate

import (
	"fmt"
	"strings"
)

const defaultMaxLen = 80

// Truncator trims long secret values to a maximum length.
type Truncator struct {
	maxLen int
	suffix string
}

// New returns a Truncator with the default maximum length.
func New() *Truncator {
	return &Truncator{maxLen: defaultMaxLen, suffix: "..."}
}

// NewWithLimit returns a Truncator with a custom maximum length.
func NewWithLimit(maxLen int) *Truncator {
	if maxLen < 1 {
		maxLen = defaultMaxLen
	}
	return &Truncator{maxLen: maxLen, suffix: "..."}
}

// Value truncates a single string value if it exceeds the max length.
func (t *Truncator) Value(s string) string {
	if len(s) <= t.maxLen {
		return s
	}
	return s[:t.maxLen] + t.suffix
}

// Map applies truncation to all values in a map, returning a new map.
func (t *Truncator) Map(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = t.Value(v)
	}
	return out
}

// ContainsLong reports whether any value in the map exceeds the max length.
func (t *Truncator) ContainsLong(secrets map[string]string) bool {
	for _, v := range secrets {
		if len(v) > t.maxLen {
			return true
		}
	}
	return false
}

// Summary returns a short description of how many values were truncated.
func (t *Truncator) Summary(secrets map[string]string) string {
	count := 0
	for _, v := range secrets {
		if len(v) > t.maxLen {
			count++
		}
	}
	if count == 0 {
		return ""
	}
	plural := "values"
	if count == 1 {
		plural = "value"
	}
	_ = strings.Join // ensure import used
	return fmt.Sprintf("%d %s truncated to %d characters", count, plural, t.maxLen)
}
