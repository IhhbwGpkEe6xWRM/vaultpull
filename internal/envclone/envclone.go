// Package envclone provides utilities for deep-copying secret maps
// with optional key filtering and value transformation hooks.
package envclone

import "fmt"

// Option configures the Cloner.
type Option func(*Cloner)

// WithKeyFilter restricts cloning to keys accepted by fn.
func WithKeyFilter(fn func(string) bool) Option {
	return func(c *Cloner) { c.keyFilter = fn }
}

// WithValueHook applies fn to every value before it is placed in the clone.
func WithValueHook(fn func(string, string) (string, error)) Option {
	return func(c *Cloner) { c.valueHook = fn }
}

// Cloner copies secret maps.
type Cloner struct {
	keyFilter func(string) bool
	valueHook func(string, string) (string, error)
}

// New returns a Cloner configured with opts.
func New(opts ...Option) *Cloner {
	c := &Cloner{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Clone returns a deep copy of src, applying any registered filter and hook.
// If a value hook returns an error the clone is aborted and the error is
// returned together with a nil map.
func (c *Cloner) Clone(src map[string]string) (map[string]string, error) {
	if src == nil {
		return nil, fmt.Errorf("envclone: source map must not be nil")
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if c.keyFilter != nil && !c.keyFilter(k) {
			continue
		}
		if c.valueHook != nil {
			var err error
			v, err = c.valueHook(k, v)
			if err != nil {
				return nil, fmt.Errorf("envclone: hook error on key %q: %w", k, err)
			}
		}
		out[k] = v
	}
	return out, nil
}
