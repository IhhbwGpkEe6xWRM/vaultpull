// Package envcast provides type-casting utilities for secret map values.
// It converts raw string values from Vault into typed Go primitives.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Caster converts string values from a secret map into typed primitives.
type Caster struct{}

// New returns a new Caster.
func New() *Caster {
	return &Caster{}
}

// String returns the raw string value for the given key.
func (c *Caster) String(m map[string]string, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("envcast: key %q not found", key)
	}
	return v, nil
}

// Int parses the value for key as an integer.
func (c *Caster) Int(m map[string]string, key string) (int64, error) {
	v, err := c.String(m, key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("envcast: key %q: cannot parse %q as int: %w", key, v, err)
	}
	return n, nil
}

// Bool parses the value for key as a boolean.
// Accepts: true, false, 1, 0, yes, no (case-insensitive).
func (c *Caster) Bool(m map[string]string, key string) (bool, error) {
	v, err := c.String(m, key)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("envcast: key %q: cannot parse %q as bool", key, v)
	}
}

// Float parses the value for key as a float64.
func (c *Caster) Float(m map[string]string, key string) (float64, error) {
	v, err := c.String(m, key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0, fmt.Errorf("envcast: key %q: cannot parse %q as float: %w", key, v, err)
	}
	return f, nil
}
