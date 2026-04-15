// Package transform provides key/value transformation utilities for
// normalising secrets before they are written to .env files.
//
// Transformations are applied in the order they are registered and are
// composable: each Transformer receives the output of the previous one.
package transform

import (
	"strings"
	"unicode"
)

// Transformer is a function that transforms a single key/value pair.
// Returning an empty key signals that the entry should be dropped.
type Transformer func(key, value string) (string, string)

// Pipeline applies a sequence of Transformers to a map of secrets,
// returning a new map with all transformations applied.
type Pipeline struct {
	steps []Transformer
}

// New creates a Pipeline with the provided Transformers applied in order.
func New(steps ...Transformer) *Pipeline {
	return &Pipeline{steps: steps}
}

// Apply runs every Transformer in the pipeline over each entry in src.
// Entries whose key becomes empty after any step are omitted from the result.
func (p *Pipeline) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		for _, step := range p.steps {
			k, v = step(k, v)
			if k == "" {
				break
			}
		}
		if k != "" {
			out[k] = v
		}
	}
	return out
}

// UppercaseKeys returns a Transformer that converts every key to upper-case.
func UppercaseKeys() Transformer {
	return func(key, value string) (string, string) {
		return strings.ToUpper(key), value
	}
}

// PrefixKeys returns a Transformer that prepends prefix to every key.
func PrefixKeys(prefix string) Transformer {
	return func(key, value string) (string, string) {
		if prefix == "" {
			return key, value
		}
		return prefix + key, value
	}
}

// TrimSpace returns a Transformer that trims leading/trailing whitespace
// from both keys and values.
func TrimSpace() Transformer {
	return func(key, value string) (string, string) {
		return strings.TrimSpace(key), strings.TrimSpace(value)
	}
}

// ReplaceHyphens returns a Transformer that replaces hyphens in keys with
// underscores, which is required for valid shell variable names.
func ReplaceHyphens() Transformer {
	return func(key, value string) (string, string) {
		return strings.ReplaceAll(key, "-", "_"), value
	}
}

// DropNonPrintable returns a Transformer that removes non-printable runes
// from values, guarding against accidental control characters in secrets.
func DropNonPrintable() Transformer {
	return func(key, value string) (string, string) {
		cleaned := strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, value)
		return key, cleaned
	}
}
