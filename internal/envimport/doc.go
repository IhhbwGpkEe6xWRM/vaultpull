// Package envimport provides utilities for reading an existing .env file
// and merging its contents with secrets fetched from Vault.
//
// Basic usage:
//
//	importer := envimport.New(".env")
//	local, err := importer.Load()
//
//	// merge with vault secrets, vault values win on conflict
//	merged, err := envimport.Merge(local, vaultSecrets, envimport.VaultWins)
//
// The Parse function accepts any io.Reader, making it easy to test without
// touching the filesystem.
package envimport
