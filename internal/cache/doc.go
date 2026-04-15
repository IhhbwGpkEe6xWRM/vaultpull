// Package cache implements a lightweight, file-backed cache for Vault secret
// snapshots used by vaultpull.
//
// The cache stores a SHA-256 checksum of each secret map alongside the raw
// key/value pairs. On subsequent runs, vaultpull can compare the current
// Vault response against the cached checksum to decide whether the local
// .env file needs to be rewritten, avoiding unnecessary disk writes and
// making audit logs easier to interpret.
//
// Cache files are written as newline-delimited JSON with file mode 0600 so
// that secret values stored in the cache are not world-readable.
//
// Usage:
//
//	store, err := cache.NewStore(".vaultpull.cache")
//	if err != nil { ... }
//
//	if store.IsFresh(path, secrets) {
//	    // skip writing .env
//	}
//	store.Set(path, secrets)
//	_ = store.Save()
package cache
