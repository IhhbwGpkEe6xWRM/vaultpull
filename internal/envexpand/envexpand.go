// Package envexpand resolves variable references within secret values.
// References of the form ${KEY} or $KEY are expanded using values from
// the same secret map or a provided fallback environment.
package envexpand

import (
	"os"
	"strings"
)

// Expander resolves variable references in secret values.
type Expander struct {
	useOS bool
}

// Option configures an Expander.
type Option func(*Expander)

// WithOSFallback allows falling back to os.Getenv when a key is not
// present in the secret map.
func WithOSFallback() Option {
	return func(e *Expander) { e.useOS = true }
}

// New returns a new Expander with the given options.
func New(opts ...Option) *Expander {
	e := &Expander{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Expand returns a new map where every value has its ${KEY} and $KEY
// references replaced. Self-references are left unexpanded to avoid
// infinite loops.
func (e *Expander) Expand(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = e.expandValue(k, v, secrets)
	}
	return out
}

func (e *Expander) expandValue(self, value string, secrets map[string]string) string {
	return os.Expand(value, func(key string) string {
		if key == self {
			return "" // prevent self-reference loop
		}
		if v, ok := secrets[key]; ok {
			return v
		}
		if e.useOS {
			return os.Getenv(key)
		}
		return ""
	})
}

// ContainsReferences reports whether the value contains any $KEY or
// ${KEY} style references.
func ContainsReferences(value string) bool {
	return strings.Contains(value, "$")
}
