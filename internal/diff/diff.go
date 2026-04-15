// Package diff provides utilities for comparing secret maps and producing
// human-readable change summaries before writing .env files.
package diff

import "sort"

// ChangeType represents the kind of change detected for a secret key.
type ChangeType int

const (
	Added ChangeType = iota
	Removed
	Modified
	Unchanged
)

// Change describes a single key-level difference between two secret maps.
type Change struct {
	Key  string
	Type ChangeType
}

// Result holds the full diff between an old and new secret map.
type Result struct {
	Changes []Change
}

// HasChanges returns true if any keys were added, removed, or modified.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns counts of each change type.
func (r *Result) Summary() (added, removed, modified int) {
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return
}

// Compare produces a Result describing the differences between old and new
// secret maps. Values in old/new are compared by equality; keys present in
// old but absent from new are Removed, and vice-versa for Added.
func Compare(old, next map[string]string) Result {
	seen := make(map[string]bool)
	var changes []Change

	for k, newVal := range next {
		seen[k] = true
		if oldVal, exists := old[k]; !exists {
			changes = append(changes, Change{Key: k, Type: Added})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: Modified})
		} else {
			changes = append(changes, Change{Key: k, Type: Unchanged})
		}
	}

	for k := range old {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Type: Removed})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Changes: changes}
}
