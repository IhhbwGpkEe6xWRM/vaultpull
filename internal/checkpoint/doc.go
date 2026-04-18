// Package checkpoint records the last successful sync state for each secret
// path, including a content checksum and timestamp.
//
// It allows vaultpull to skip re-writing .env files when the remote secrets
// have not changed since the previous sync, reducing unnecessary disk writes
// and downstream process restarts.
//
// Usage:
//
//	store, err := checkpoint.NewStore(".vaultpull/checkpoint.json")
//	if err != nil { ... }
//
//	if store.IsFresh(path, checksum, 5*time.Minute, time.Now()) {
//		// skip sync
//	}
//
//	store.Set(path, checksum, time.Now())
package checkpoint
