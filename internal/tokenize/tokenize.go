package tokenize

import (
	"fmt"
	"regexp"
	"strings"
)

// Token represents a single parsed segment from a secret path.
type Token struct {
	Raw   string
	Parts []string
}

// Tokenizer splits secret paths into structured tokens for matching or display.
type Tokenizer struct {
	separator string
	pattern   *regexp.Regexp
}

// New returns a Tokenizer using the default "/" separator.
func New() *Tokenizer {
	return NewWithSeparator("/")
}

// NewWithSeparator returns a Tokenizer with a custom separator.
func NewWithSeparator(sep string) *Tokenizer {
	escaped := regexp.QuoteMeta(sep)
	return &Tokenizer{
		separator: sep,
		pattern:   regexp.MustCompile(escaped + "+"),
	}
}

// Parse splits a raw path string into a Token.
func (t *Tokenizer) Parse(raw string) Token {
	trimmed := strings.Trim(raw, t.separator)
	if trimmed == "" {
		return Token{Raw: raw, Parts: []string{}}
	}
	parts := t.pattern.Split(trimmed, -1)
	filtered := parts[:0]
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return Token{Raw: raw, Parts: filtered}
}

// Join reconstructs a path from a Token's parts.
func (t *Tokenizer) Join(tok Token) string {
	return strings.Join(tok.Parts, t.separator)
}

// Depth returns the number of segments in the token.
func (t *Tokenizer) Depth(tok Token) int {
	return len(tok.Parts)
}

// Parent returns a Token representing the parent path, or an error if already root.
func (t *Tokenizer) Parent(tok Token) (Token, error) {
	if len(tok.Parts) == 0 {
		return Token{}, fmt.Errorf("tokenize: path %q has no parent", tok.Raw)
	}
	parentParts := tok.Parts[:len(tok.Parts)-1]
	parentRaw := strings.Join(parentParts, t.separator)
	return Token{Raw: parentRaw, Parts: parentParts}, nil
}
