// Package envguard detects write conflicts between locally modified .env files
// and incoming secrets fetched from Vault.
//
// A conflict occurs when both the local file and the remote Vault secret have
// diverged from the last-synced snapshot for the same key. This allows
// vaultpull to warn users before silently overwriting hand-edited values.
//
// Basic usage:
//
//	g := envguard.New(lastSyncedSnapshot)
//	violations, err := g.Check(localEnvMap, incomingVaultMap)
//	if len(violations) > 0 {
//		fmt.Print(envguard.Summary(violations))
//	}
package envguard
