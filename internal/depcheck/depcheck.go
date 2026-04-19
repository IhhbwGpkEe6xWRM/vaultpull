// Package depcheck verifies that required secrets exist before writing
// the output .env file, preventing partial or broken configurations.
package depcheck

import "fmt"

// Violation describes a single missing or empty required secret.
type Violation struct {
	Key    string
	Reason string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Reason)
}

// Checker holds a set of required keys.
type Checker struct {
	required []string
}

// New returns a Checker that enforces the given required keys.
func New(required []string) *Checker {
	keys := make([]string, len(required))
	copy(keys, required)
	return &Checker{required: keys}
}

// Check returns a list of violations for any required key that is
// absent from secrets or whose value is an empty string.
func (c *Checker) Check(secrets map[string]string) []Violation {
	var violations []Violation
	for _, k := range c.required {
		v, ok := secrets[k]
		if !ok {
			violations = append(violations, Violation{Key: k, Reason: "missing"})
			continue
		}
		if v == "" {
			violations = append(violations, Violation{Key: k, Reason: "empty value"})
		}
	}
	return violations
}

// Summary returns a human-readable summary of violations, or an empty
// string when there are none.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return ""
	}
	msg := fmt.Sprintf("%d required secret(s) failed checks:\n", len(violations))
	for _, v := range violations {
		msg += fmt.Sprintf("  - %s\n", v.Error())
	}
	return msg
}
