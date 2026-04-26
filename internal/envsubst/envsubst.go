// Package envsubst performs shell-style variable substitution on env map values.
// References of the form ${KEY} or $KEY are replaced with their resolved values.
package envsubst

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Substitutor replaces variable references in env values.
type Substitutor struct {
	maxDepth int
	os       map[string]string
}

// Option configures a Substitutor.
type Option func(*Substitutor)

// WithMaxDepth sets the maximum substitution pass depth (default 5).
func WithMaxDepth(d int) Option {
	return func(s *Substitutor) {
		if d > 0 {
			s.maxDepth = d
		}
	}
}

// WithOSEnv provides fallback values from the OS environment.
func WithOSEnv(env map[string]string) Option {
	return func(s *Substitutor) {
		s.os = env
	}
}

// New creates a Substitutor with the given options.
func New(opts ...Option) *Substitutor {
	s := &Substitutor{maxDepth: 5}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply performs variable substitution on all values in src.
// Values in src take precedence over OS fallback values.
func (s *Substitutor) Apply(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	for range make([]struct{}, s.maxDepth) {
		changed := false
		for k, v := range out {
			replaced, err := s.expand(v, out)
			if err != nil {
				return nil, fmt.Errorf("envsubst: key %q: %w", k, err)
			}
			if replaced != v {
				out[k] = replaced
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return out, nil
}

func (s *Substitutor) expand(val string, src map[string]string) (string, error) {
	return refPattern.ReplaceAllStringFunc(val, func(match string) string {
		key := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		if v, ok := src[key]; ok {
			return v
		}
		if s.os != nil {
			if v, ok := s.os[key]; ok {
				return v
			}
		}
		return ""
	}), nil
}

// ContainsReferences reports whether the value contains any substitution references.
func ContainsReferences(val string) bool {
	return refPattern.MatchString(val)
}
