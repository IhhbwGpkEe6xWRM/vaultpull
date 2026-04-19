package envtag

// Filter selects entries from a slice of tagged secret sets based on
// required tags. Each entry is a map of secrets with an associated Tagger.
type Entry struct {
	Path    string
	Secrets map[string]string
	Tagger  *Tagger
}

// Filter returns only entries whose tagger matches all required filter tags.
func Filter(entries []Entry, filter []Tag) []Entry {
	if len(filter) == 0 {
		return entries
	}
	var out []Entry
	for _, e := range entries {
		if e.Tagger != nil && e.Tagger.MatchesAll(filter) {
			out = append(out, e)
		}
	}
	return out
}

// ParseFilter parses a slice of "key:value" strings into a Tag slice
// suitable for passing to Filter or MatchesAll.
func ParseFilter(raw []string) []Tag {
	tagger := New(raw)
	return tagger.Tags()
}
