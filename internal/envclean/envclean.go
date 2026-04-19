// Package envclean removes stale keys from a local .env file that no longer
// exist in the Vault secret map.
package envclean

import "sort"

// Result holds the outcome of a clean operation.
type Result struct {
	Removed []string
	Kept    map[string]string
}

// Cleaner removes keys from local that are absent in incoming.
type Cleaner struct {
	dryRun bool
}

// Option configures a Cleaner.
type Option func(*Cleaner)

// WithDryRun makes Clean report removals without modifying the map.
func WithDryRun() Option {
	return func(c *Cleaner) { c.dryRun = true }
}

// New returns a new Cleaner.
func New(opts ...Option) *Cleaner {
	c := &Cleaner{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Clean returns a Result where keys present in local but absent in incoming
// are listed as Removed. If dry-run is disabled the Kept map omits those keys.
func (c *Cleaner) Clean(local, incoming map[string]string) Result {
	removed := []string{}
	kept := make(map[string]string, len(local))

	for k, v := range local {
		if _, ok := incoming[k]; !ok {
			removed = append(removed, k)
		} else {
			kept[k] = v
		}
	}

	sort.Strings(removed)

	if c.dryRun {
		// In dry-run mode return original local so nothing is mutated.
		full := make(map[string]string, len(local))
		for k, v := range local {
			full[k] = v
		}
		return Result{Removed: removed, Kept: full}
	}

	return Result{Removed: removed, Kept: kept}
}
