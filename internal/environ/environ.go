// Package environ provides utilities for reading and merging environment
// variables with secret maps, allowing local overrides to take precedence.
package environ

import (
	"os"
	"strings"
)

// Loader reads environment variables and merges them with a secrets map.
type Loader struct {
	prefix  string
	override bool
}

// New returns a Loader that matches variables with the given prefix.
// If override is true, existing env vars take precedence over vault secrets.
func New(prefix string, override bool) *Loader {
	return &Loader{
		prefix:   strings.ToUpper(strings.Trim(prefix, "_")),
		override: override,
	}
}

// Load returns all environment variables matching the prefix as a map.
// Keys are returned without the prefix.
func (l *Loader) Load() map[string]string {
	result := make(map[string]string)
	for _, kv := range os.Environ() {
		key, val, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		upper := strings.ToUpper(key)
		if l.prefix == "" || strings.HasPrefix(upper, l.prefix+"_") {
			stripped := key
			if l.prefix != "" {
				stripped = key[len(l.prefix)+1:]
			}
			result[stripped] = val
		}
	}
	return result
}

// Merge combines env vars with the provided secrets map.
// When override is true, env vars win on conflict; otherwise secrets win.
func (l *Loader) Merge(secrets map[string]string) map[string]string {
	env := l.Load()
	out := make(map[string]string, len(secrets)+len(env))
	for k, v := range secrets {
		out[k] = v
	}
	for k, v := range env {
		if _, exists := out[k]; !exists || l.override {
			out[k] = v
		}
	}
	return out
}
