// Package schema provides validation of secret maps against a declared
// schema, ensuring required keys are present and values match expected
// types or patterns.
package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldType represents the expected type of a secret value.
type FieldType string

const (
	TypeString FieldType = "string"
	TypeInt    FieldType = "int"
	TypeBool   FieldType = "bool"
)

// Field describes a single expected secret key.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // optional regex
}

// Schema holds a collection of field definitions.
type Schema struct {
	fields []Field
}

// New creates a Schema from the given field definitions.
func New(fields []Field) *Schema {
	return &Schema{fields: fields}
}

// Violation describes a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks secrets against the schema and returns all violations.
func (s *Schema) Validate(secrets map[string]string) []Violation {
	var violations []Violation

	for _, f := range s.fields {
		val, ok := secrets[f.Key]
		if !ok || val == "" {
			if f.Required {
				violations = append(violations, Violation{Key: f.Key, Message: "required key is missing or empty"})
			}
			continue
		}

		switch f.Type {
		case TypeInt:
			if !isInt(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("expected int, got %q", val)})
			}
		case TypeBool:
			if !isBool(val) {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("expected bool, got %q", val)})
			}
		}

		if f.Pattern != "" {
			matched, err := regexp.MatchString(f.Pattern, val)
			if err != nil || !matched {
				violations = append(violations, Violation{Key: f.Key, Message: fmt.Sprintf("value does not match pattern %q", f.Pattern)})
			}
		}
	}

	return violations
}

func isInt(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func isBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "false", "1", "0", "yes", "no":
		return true
	}
	return false
}
