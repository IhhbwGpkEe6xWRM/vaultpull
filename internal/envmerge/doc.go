// Package envmerge provides utilities for merging Vault secrets with existing
// local .env key-value maps.
//
// Three strategies are supported:
//
//   - VaultWins: Vault values overwrite local values on collision (default).
//   - LocalWins: Local values are preserved when a key exists in both sources.
//   - ErrorOnConflict: Returns an error listing all conflicting keys so the
//     caller can decide how to proceed.
//
// Usage:
//
//	result, err := envmerge.Merge(localMap, vaultMap, envmerge.VaultWins)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Overridden) // keys that were replaced
package envmerge
