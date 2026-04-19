package envimport

import "fmt"

// ConflictError describes a key present in both local and vault with different values.
type ConflictError struct {
	Key        string
	LocalValue string
	VaultValue string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("envimport: conflict on key %q: local=%q vault=%q",
		e.Key, e.LocalValue, e.VaultValue)
}

// MergeStrategy controls how conflicts are resolved.
type MergeStrategy int

const (
	// VaultWins overwrites local values with vault values.
	VaultWins MergeStrategy = iota
	// LocalWins keeps local values when a conflict exists.
	LocalWins
	// ErrorOnConflict returns an error when the same key has different values.
	ErrorOnConflict
)

// Merge combines local (from .env file) and vault secrets using the given strategy.
// Keys present only in one source are always included.
func Merge(local, vault map[string]string, strategy MergeStrategy) (map[string]string, error) {
	out := make(map[string]string, len(local)+len(vault))
	for k, v := range local {
		out[k] = v
	}
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
		case LocalWins:
			// keep existing
		case ErrorOnConflict:
			return nil, &ConflictError{Key: k, LocalValue: localVal, VaultValue: vaultVal}
		}
	}
	return out, nil
}
