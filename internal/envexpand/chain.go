package envexpand

// MultiPass runs Expand up to maxPasses times so that chained references
// (e.g. A=${B}, B=${C}, C=hello) are fully resolved. Expansion stops early
// when a pass produces no further changes.
func (e *Expander) MultiPass(secrets map[string]string, maxPasses int) map[string]string {
	if maxPasses <= 0 {
		maxPasses = 5
	}
	current := copyMap(secrets)
	for i := 0; i < maxPasses; i++ {
		next := e.Expand(current)
		if mapsEqual(current, next) {
			break
		}
		current = next
	}
	return current
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
