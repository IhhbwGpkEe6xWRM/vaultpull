package envscope

import (
	"fmt"
	"sort"
)

// Guard wraps a Scope and tracks write attempts, blocking any key that falls
// outside the declared scope and recording which keys were blocked.
type Guard struct {
	scope   *Scope
	blocked []string
}

// NewGuard creates a Guard backed by the given Scope.
func NewGuard(s *Scope) *Guard {
	return &Guard{scope: s}
}

// Write attempts to add key/value to dst. If the key is outside the scope it
// is recorded as blocked and an error is returned; dst is not modified.
func (g *Guard) Write(dst map[string]string, key, value string) error {
	if !g.scope.Allows(key) {
		g.blocked = append(g.blocked, key)
		return fmt.Errorf("envscope: write blocked for key %q: outside scope", key)
	}
	dst[key] = value
	return nil
}

// WriteAll applies Write for every entry in src. It continues on error and
// returns a combined error listing all blocked keys at the end.
func (g *Guard) WriteAll(dst, src map[string]string) error {
	for k, v := range src {
		_ = g.Write(dst, k, v) // errors collected via g.blocked
	}
	if len(g.blocked) == 0 {
		return nil
	}
	sorted := make([]string, len(g.blocked))
	copy(sorted, g.blocked)
	sort.Strings(sorted)
	return fmt.Errorf("envscope: %d key(s) blocked: %v", len(sorted), sorted)
}

// Blocked returns a sorted copy of all keys that were blocked so far.
func (g *Guard) Blocked() []string {
	out := make([]string, len(g.blocked))
	copy(out, g.blocked)
	sort.Strings(out)
	return out
}

// Reset clears the blocked-key history.
func (g *Guard) Reset() {
	g.blocked = g.blocked[:0]
}
