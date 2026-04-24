// Package envpatch applies a partial update (patch) to an existing secrets map,
// merging only the changed or new keys while optionally removing deleted ones.
package envpatch

import "fmt"

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
)

// Entry describes a single change in a patch.
type Entry struct {
	Key   string
	Value string
	Op    Op
}

// Patch holds a list of entries to apply.
type Patch struct {
	entries []Entry
}

// New creates a Patch from a slice of entries. Returns an error if any entry
// has an empty key or an unrecognised operation.
func New(entries []Entry) (*Patch, error) {
	for i, e := range entries {
		if e.Key == "" {
			return nil, fmt.Errorf("envpatch: entry %d has empty key", i)
		}
		if e.Op != OpSet && e.Op != OpDelete {
			return nil, fmt.Errorf("envpatch: entry %d has unknown op %q", i, e.Op)
		}
	}
	return &Patch{entries: entries}, nil
}

// Apply merges the patch into base and returns a new map. The original map is
// never mutated.
func (p *Patch) Apply(base map[string]string) map[string]string {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for _, e := range p.entries {
		switch e.Op {
		case OpSet:
			out[e.Key] = e.Value
		case OpDelete:
			delete(out, e.Key)
		}
	}
	return out
}

// Diff returns the entries that differ between base and patched, expressed as
// patch entries. Keys present only in patched become OpSet; keys missing from
// patched that were in base become OpDelete.
func Diff(base, patched map[string]string) []Entry {
	var entries []Entry
	for k, v := range patched {
		if bv, ok := base[k]; !ok || bv != v {
			entries = append(entries, Entry{Key: k, Value: v, Op: OpSet})
		}
	}
	for k := range base {
		if _, ok := patched[k]; !ok {
			entries = append(entries, Entry{Key: k, Op: OpDelete})
		}
	}
	return entries
}
