// Package envguard provides protection against overwriting env file keys
// that have been locally modified since the last sync.
package envguard

import (
	"fmt"
	"sort"
)

// Violation describes a key that would be overwritten.
type Violation struct {
	Key      string
	Local    string
	Incoming string
}

// Guard checks incoming secrets against a locally parsed env map and
// returns any keys whose local value differs from the last-synced snapshot.
type Guard struct {
	protected map[string]string // last-synced snapshot
}

// New creates a Guard using the provided last-synced snapshot.
func New(snapshot map[string]string) *Guard {
	cp := make(map[string]string, len(snapshot))
	for k, v := range snapshot {
		cp[k] = v
	}
	return &Guard{protected: cp}
}

// Check compares local (current on-disk) values against incoming (from Vault)
// values. A violation is raised when:
//   - the key exists in local and incoming
//   - the local value differs from the last-synced snapshot value
//   - the incoming value differs from the last-synced snapshot value
//
// i.e. both sides changed — a conflict.
func (g *Guard) Check(local, incoming map[string]string) ([]Violation, error) {
	var violations []Violation

	for key, inVal := range incoming {
		snapshotVal, wasSynced := g.protected[key]
		if !wasSynced {
			continue
		}
		localVal, existsLocally := local[key]
		if !existsLocally {
			continue
		}
		localDrifted := localVal != snapshotVal
		incomingChanged := inVal != snapshotVal
		if localDrifted && incomingChanged {
			violations = append(violations, Violation{
				Key:      key,
				Local:    localVal,
				Incoming: inVal,
			})
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})
	return violations, nil
}

// Summary returns a human-readable summary of violations.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return "no conflicts detected"
	}
	msg := fmt.Sprintf("%d conflict(s) detected:\n", len(violations))
	for _, v := range violations {
		msg += fmt.Sprintf("  %s: local=%q incoming=%q\n", v.Key, v.Local, v.Incoming)
	}
	return msg
}
