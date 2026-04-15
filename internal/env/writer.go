// Package env provides utilities for writing secrets to .env files.
package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing key-value pairs to a .env file.
type Writer struct {
	path string
}

// NewWriter creates a new Writer targeting the given file path.
func NewWriter(path string) *Writer {
	return &Writer{path: path}
}

// Write serializes the provided secrets map into .env format and writes it
// to the configured file path, truncating any existing content.
func (w *Writer) Write(secrets map[string]string) error {
	f, err := os.Create(w.path)
	if err != nil {
		return fmt.Errorf("env: creating file %q: %w", w.path, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", sanitizeKey(k), quoteValue(secrets[k]))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env: writing key %q: %w", k, err)
		}
	}

	return nil
}

// sanitizeKey uppercases the key and replaces hyphens and spaces with underscores
// to produce a valid environment variable name.
func sanitizeKey(k string) string {
	k = strings.ToUpper(k)
	k = strings.ReplaceAll(k, "-", "_")
	k = strings.ReplaceAll(k, " ", "_")
	return k
}

// quoteValue wraps value in double quotes if it contains whitespace or special
// characters that could break shell parsing.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r#$\"\'\\`") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
