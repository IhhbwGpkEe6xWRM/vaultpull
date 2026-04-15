// Package validate provides helpers for validating secret maps
// before they are written to .env files.
package validate

import (
	"errors"
	"fmt"
	"strings"
)

// Issue describes a single validation problem found in a secret map.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("key %q: %s", i.Key, i.Message)
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// OK returns true when no issues were found.
func (r Result) OK() bool { return len(r.Issues) == 0 }

// Err returns a combined error when issues exist, otherwise nil.
func (r Result) Err() error {
	if r.OK() {
		return nil
	}
	msgs := make([]string, len(r.Issues))
	for i, iss := range r.Issues {
		msgs[i] = iss.Error()
	}
	return errors.New(strings.Join(msgs, "; "))
}

// Secrets validates a map of secret key/value pairs and returns a Result.
// Rules:
//   - Keys must not be empty.
//   - Keys must contain only ASCII letters, digits, and underscores.
//   - Values must not exceed maxValueLen bytes.
const maxValueLen = 65536

func Secrets(secrets map[string]string) Result {
	var result Result
	for k, v := range secrets {
		if k == "" {
			result.Issues = append(result.Issues, Issue{Key: k, Message: "key must not be empty"})
			continue
		}
		if !isValidKey(k) {
			result.Issues = append(result.Issues, Issue{Key: k, Message: "key contains invalid characters (allowed: A-Z, a-z, 0-9, _)"})
		}
		if len(v) > maxValueLen {
			result.Issues = append(result.Issues, Issue{Key: k, Message: fmt.Sprintf("value exceeds maximum length of %d bytes", maxValueLen)})
		}
	}
	return result
}

func isValidKey(s string) bool {
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
