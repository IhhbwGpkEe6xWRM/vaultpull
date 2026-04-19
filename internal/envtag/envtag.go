// Package envtag provides tagging and filtering of secret keys using
// structured tag annotations (e.g. "env:production", "team:platform").
package envtag

import (
	"strings"
)

// Tag represents a single key:value annotation.
type Tag struct {
	Key   string
	Value string
}

// Tagger holds a set of tags applied to a secret map.
type Tagger struct {
	tags []Tag
}

// New parses raw tag strings of the form "key:value" and returns a Tagger.
// Malformed entries (missing colon or empty key) are silently ignored.
func New(raw []string) *Tagger {
	var tags []Tag
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			continue
		}
		tags = append(tags, Tag{Key: k, Value: v})
	}
	return &Tagger{tags: tags}
}

// Tags returns the parsed tags.
func (t *Tagger) Tags() []Tag {
	return t.tags
}

// Annotate adds tag metadata as special keys into the secrets map using the
// prefix "__tag_<key>" so downstream tooling can inspect provenance.
func (t *Tagger) Annotate(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets)+len(t.tags))
	for k, v := range secrets {
		out[k] = v
	}
	for _, tag := range t.tags {
		key := "__tag_" + strings.ToUpper(tag.Key)
		out[key] = tag.Value
	}
	return out
}

// MatchesAll returns true if the provided filter tags are all present in t.
func (t *Tagger) MatchesAll(filter []Tag) bool {
	index := make(map[string]string, len(t.tags))
	for _, tag := range t.tags {
		index[tag.Key] = tag.Value
	}
	for _, f := range filter {
		if v, ok := index[f.Key]; !ok || v != f.Value {
			return false
		}
	}
	return true
}
