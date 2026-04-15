// Package snapshot captures and persists the state of secrets fetched from
// HashiCorp Vault, enabling vaultpull to detect changes between runs.
//
// # Overview
//
// A Snapshot records the full set of key-value secrets for a given Vault path
// and namespace at a specific point in time. Snapshots are stored as JSON
// files on disk with restricted permissions (0600).
//
// # Usage
//
// Create a Store pointing at a snapshot file, then Save or Load snapshots:
//
//	store := snapshot.NewStore(".vaultpull.snapshot.json")
//
//	prev, err := store.Load()   // nil if first run
//	changes, err := snapshot.Report(os.Stdout, prev, currentSecrets)
//
//	err = store.Save(snapshot.Snapshot{
//		Path:    "secret/myapp",
//		Secrets: currentSecrets,
//	})
//
// The Report function uses the diff package internally and writes a
// human-readable summary of added, removed, and modified keys.
package snapshot
