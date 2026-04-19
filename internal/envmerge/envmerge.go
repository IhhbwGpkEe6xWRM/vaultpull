// Package envmerge combines secrets from Vault with existing local .env values,
// applying a configurable precedence strategy.
package envmerge

import (
	"fmt"
	"sort"
	"strings"
)

// Strategy controls which source wins on key collision.
type Strategy int

const (
	// VaultWins overwrites local values with Vault secrets (default).
	VaultWins Strategy = iota
	// LocalWins keeps the local value when a key exists in both sources.
	LocalWins
	// ErrorOnConflict returns an error if the same key appears in both sources with different values.
	ErrorOnConflict
)

// Result holds the merged map and metadata about the merge.
type Result struct {
	Secrets    map[string]string
	Overridden []string // keys where the losing source was overwritten
	Conflicts  []string // keys in conflict (only populated for ErrorOnConflict)
}

// Merge combines local and vault maps according to the given strategy.
func Merge(local, vault map[string]string, strategy Strategy) (*Result, error) {
	out := make(map[string]string, len(local)+len(vault))
	for k, v := range local {
		out[k] = v
	}

	var overridden []string
	var conflicts []string

	for k, vaultVal := range vault {
		localVal, exists := out[k]
		if !exists {
			out[k] = vaultVal
			continue
		}
		if localVal == vaultVal {
			continue
		}
		switch strategy {
		case VaultWins:
			out[k] = vaultVal
			overridden = append(overridden, k)
		case LocalWins:
			// keep local; note it as overridden from vault perspective
			overridden = append(overridden, k)
		case ErrorOnConflict:
			conflicts = append(conflicts, k)
		}
	}

	if len(conflicts) > 0 {
		sort.Strings(conflicts)
		return nil, fmt.Errorf("envmerge: conflicts on keys: %s", strings.Join(conflicts, ", "))
	}

	sort.Strings(overridden)
	return &Result{Secrets: out, Overridden: overridden}, nil
}
