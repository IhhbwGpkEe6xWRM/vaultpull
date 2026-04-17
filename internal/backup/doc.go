// Package backup provides pre-write backup functionality for .env files.
//
// Before vaultpull overwrites a local .env file it can save a timestamped
// copy to a configurable backup directory.  Backups can later be listed and
// restored if a sync introduces unintended changes.
//
// Usage:
//
//	store, err := backup.NewStore(".vaultpull/backups")
//	if err != nil { ... }
//
//	// Save current .env before overwriting
//	path, err := store.Save(".env")
//
//	// List all backups for .env
//	files, err := store.List(".env")
//
//	// Restore a specific backup
//	err = store.Restore(files[0], ".env")
package backup
