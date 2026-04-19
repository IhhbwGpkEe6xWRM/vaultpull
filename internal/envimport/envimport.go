// Package envimport parses an existing .env file into a secrets map.
package envimport

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Importer reads a .env file and returns key/value pairs.
type Importer struct {
	path string
}

// New returns an Importer that reads from path.
func New(path string) *Importer {
	return &Importer{path: path}
}

// Load opens the file and parses it into a map.
func (i *Importer) Load() (map[string]string, error) {
	f, err := os.Open(i.path)
	if err != nil {
		return nil, fmt.Errorf("envimport: open %q: %w", i.path, err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse reads .env lines from r and returns a map.
// Lines starting with # are treated as comments.
// Inline comments (value # comment) are not stripped.
func Parse(r io.Reader) (map[string]string, error) {
	out := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("envimport: line %d: missing '='", lineNo)
		}
		key := strings.TrimSpace(line[:idx])
		val := unquote(strings.TrimSpace(line[idx+1:]))
		if key == "" {
			return nil, fmt.Errorf("envimport: line %d: empty key", lineNo)
		}
		out[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envimport: scan: %w", err)
	}
	return out, nil
}

// unquote strips a single layer of matching quotes if present.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
