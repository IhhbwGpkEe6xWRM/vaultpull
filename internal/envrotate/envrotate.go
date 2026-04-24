// Package envrotate provides utilities for rotating secrets in a local
// .env file by replacing old values with new ones fetched from Vault,
// while preserving keys that have not changed.
package envrotate

import (
	"errors"
	"fmt"
)

// RotateFunc is a function that receives an old value and returns a new one.
type RotateFunc func(key, oldValue string) (string, error)

// Result describes the outcome of rotating a single key.
type Result struct {
	Key     string
	OldHash string
	NewHash string
	Rotated bool
}

// Rotator applies a RotateFunc to a set of secrets and returns the updated map
// together with a per-key result slice.
type Rotator struct {
	fn RotateFunc
}

// New creates a Rotator that delegates value replacement to fn.
func New(fn RotateFunc) (*Rotator, error) {
	if fn == nil {
		return nil, errors.New("envrotate: rotate func must not be nil")
	}
	return &Rotator{fn: fn}, nil
}

// Apply iterates over secrets, calls the RotateFunc for every key and returns
// the updated map along with per-key Results. If any call returns an error the
// whole operation is aborted and the error is wrapped with the offending key.
func (r *Rotator) Apply(secrets map[string]string) (map[string]string, []Result, error) {
	if secrets == nil {
		return nil, nil, errors.New("envrotate: secrets map must not be nil")
	}

	out := make(map[string]string, len(secrets))
	results := make([]Result, 0, len(secrets))

	for k, oldVal := range secrets {
		newVal, err := r.fn(k, oldVal)
		if err != nil {
			return nil, nil, fmt.Errorf("envrotate: key %q: %w", k, err)
		}
		out[k] = newVal
		results = append(results, Result{
			Key:     k,
			OldHash: shortHash(oldVal),
			NewHash: shortHash(newVal),
			Rotated: newVal != oldVal,
		})
	}
	return out, results, nil
}

// shortHash returns the first 8 characters of a simple djb2-style hash of s
// so that values are never logged in plain text.
func shortHash(s string) string {
	h := uint32(5381)
	for i := 0; i < len(s); i++ {
		h = (h << 5) + h + uint32(s[i])
	}
	return fmt.Sprintf("%08x", h)
}
