package snapshot

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourusername/vaultpull/internal/diff"
)

// ChangeKind describes the type of change detected for a secret key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single detected change between two snapshots.
type Change struct {
	Key  string
	Kind ChangeKind
}

// Report generates a human-readable change report written to w.
// It compares prev (may be nil for first run) against current secrets.
func Report(w io.Writer, prev *Snapshot, current map[string]string) ([]Change, error) {
	var oldSecrets map[string]string
	if prev != nil {
		oldSecrets = prev.Secrets
	} else {
		oldSecrets = map[string]string{}
	}

	results := diff.Compare(oldSecrets, current)
	changes := make([]Change, 0, len(results))

	for _, r := range results {
		var kind ChangeKind
		switch r.Status {
		case diff.Added:
			kind = Added
		case diff.Removed:
			kind = Removed
		case diff.Modified:
			kind = Modified
		default:
			kind = Unchanged
		}
		changes = append(changes, Change{Key: r.Key, Kind: kind})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		if c.Kind == Unchanged {
			continue
		}
		_, err := fmt.Fprintf(w, "  [%s] %s\n", c.Kind, c.Key)
		if err != nil {
			return nil, err
		}
	}
	return changes, nil
}
