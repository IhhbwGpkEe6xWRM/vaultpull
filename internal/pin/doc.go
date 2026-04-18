// Package pin implements version pinning for Vault secret paths.
//
// A pin locks a specific secret path to a known KV version,
// so that vaultpull will not advance past that version during
// a sync operation. This is useful for controlled rollouts or
// audited environments where secret changes must be approved
// before being applied locally.
//
// Pins are stored in a JSON file (default: .vaultpull-pins.json)
// and are keyed by the full secret path. Each entry records the
// target version, the time it was pinned, and an optional actor.
//
// Usage:
//
//	store, err := pin.NewStore(".vaultpull-pins.json")
//	store.Set("secret/myapp/prod", 4, "ops-team")
//	if entry, ok := store.Get("secret/myapp/prod"); ok {
//		// use entry.Version when reading from Vault
//	}
package pin
