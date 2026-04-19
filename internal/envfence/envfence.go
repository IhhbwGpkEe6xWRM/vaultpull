// Package envfence restricts which secret keys may be written to a .env file
// based on an allowlist or denylist of key patterns.
package envfence

import (
	"fmt"
	"regexp"
	"strings"
)

// Mode controls whether the pattern list is an allowlist or denylist.
type Mode int

const (
	Allow Mode = iota // only matched keys pass
	Deny              // matched keys are blocked
)

// Fence filters a map of secrets according to compiled patterns.
type Fence struct {
	mode     Mode
	patterns []*regexp.Regexp
}

// New creates a Fence from a slice of glob-style patterns (anchored regex).
// Patterns are matched case-insensitively against the key name.
func New(mode Mode, patterns []string) (*Fence, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if strings.TrimSpace(p) == "" {
			continue
		}
		re, err := regexp.Compile("(?i)^" + p + "$")
		if err != nil {
			return nil, fmt.Errorf("envfence: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Fence{mode: mode, patterns: compiled}, nil
}

// Apply returns a filtered copy of secrets according to the fence rules.
// If no patterns are configured, all keys pass regardless of mode.
func (f *Fence) Apply(secrets map[string]string) map[string]string {
	if len(f.patterns) == 0 {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string)
	for k, v := range secrets {
		matched := f.matches(k)
		if f.mode == Allow && matched {
			out[k] = v
		} else if f.mode == Deny && !matched {
			out[k] = v
		}
	}
	return out
}

// Blocked returns the keys that would be removed by Apply.
func (f *Fence) Blocked(secrets map[string]string) []string {
	var blocked []string
	for k := range secrets {
		matched := f.matches(k)
		if (f.mode == Allow && !matched) || (f.mode == Deny && matched) {
			blocked = append(blocked, k)
		}
	}
	return blocked
}

func (f *Fence) matches(key string) bool {
	for _, re := range f.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
