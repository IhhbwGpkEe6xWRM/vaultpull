// Package envfreeze provides a mechanism to lock a set of environment
// variables so that any attempt to modify or delete them is rejected.
// This is useful when certain secrets must remain stable throughout the
// lifecycle of a sync operation (e.g. credentials used by hooks).
package envfreeze

import (
	"errors"
	"fmt"
	"sort"
)

// ErrFrozen is returned when a write targets a frozen key.
var ErrFrozen = errors.New("envfreeze: key is frozen")

// Freezer holds a snapshot of frozen keys and enforces immutability.
type Freezer struct {
	frozen map[string]string
}

// New creates a Freezer from the provided map. The keys and values are
// copied so that later mutations to src do not affect the frozen set.
func New(src map[string]string) *Freezer {
	f := &Freezer{frozen: make(map[string]string, len(src))}
	for k, v := range src {
		f.frozen[k] = v
	}
	return f
}

// Apply merges incoming into base, returning an error if any key in
// incoming would change a frozen value. Non-frozen keys are merged
// without restriction.
func (f *Freezer) Apply(base, incoming map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range incoming {
		if frozen, ok := f.frozen[k]; ok && v != frozen {
			return nil, fmt.Errorf("%w: %q", ErrFrozen, k)
		}
		out[k] = v
	}
	return out, nil
}

// Keys returns the sorted list of frozen key names.
func (f *Freezer) Keys() []string {
	keys := make([]string, 0, len(f.frozen))
	for k := range f.frozen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// IsFrozen reports whether the given key is part of the frozen set.
func (f *Freezer) IsFrozen(key string) bool {
	_, ok := f.frozen[key]
	return ok
}
