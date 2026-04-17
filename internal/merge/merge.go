// Package merge provides utilities for merging secret maps from multiple
// sources with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how conflicts are resolved when the same key exists
// in more than one source map.
type Strategy int

const (
	// LastWins uses the value from the last source that defines the key.
	LastWins Strategy = iota
	// FirstWins keeps the value from the first source that defines the key.
	FirstWins
	// ErrorOnConflict returns an error if the same key appears in multiple sources.
	ErrorOnConflict
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	First  string
	Second string
}

// Result holds the merged map and any conflicts that were detected.
type Result struct {
	Secrets   map[string]string
	Conflicts []Conflict
}

// Merge combines the provided maps according to the given strategy.
// Sources are processed in order; index 0 is considered the "first" source.
func Merge(strategy Strategy, sources ...map[string]string) (Result, error) {
	out := make(map[string]string)
	var conflicts []Conflict

	for _, src := range sources {
		for k, v := range src {
			existing, exists := out[k]
			if !exists {
				out[k] = v
				continue
			}
			conflicts = append(conflicts, Conflict{Key: k, First: existing, Second: v})
			switch strategy {
			case LastWins:
				out[k] = v
			case FirstWins:
				// keep existing — no-op
			case ErrorOnConflict:
				return Result{}, fmt.Errorf("merge: conflict on key %q", k)
			}
		}
	}

	return Result{Secrets: out, Conflicts: conflicts}, nil
}
