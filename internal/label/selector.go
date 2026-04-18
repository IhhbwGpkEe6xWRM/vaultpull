package label

// Selector filters a collection of labeled items.
type Selector struct {
	filter Set
}

// NewSelector creates a Selector from a filter Set.
func NewSelector(filter Set) *Selector {
	return &Selector{filter: filter}
}

// Filter returns only the keys from secrets whose metadata satisfies the
// selector. metaFn provides the label Set for a given secret key.
func (s *Selector) Filter(secrets map[string]string, metaFn func(key string) Set) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		if Matches(metaFn(k), s.filter) {
			out[k] = v
		}
	}
	return out
}

// ParseFilter converts a slice of "key=value" strings into a Set suitable for
// use with NewSelector.
func ParseFilter(pairs []string) Set {
	return New(pairs).labels
}
