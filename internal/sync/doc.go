// Package sync provides the top-level orchestration logic for vaultpull.
//
// It wires together the Vault client (internal/vault) and the env file writer
// (internal/env) using the loaded configuration (internal/config).
//
// Typical usage:
//
//	cfg, err := config.Load()
//	if err != nil { ... }
//
//	syncer, err := sync.New(cfg)
//	if err != nil { ... }
//
//	if err := syncer.Run(); err != nil { ... }
//
// For testing, use NewWithDeps to inject mock implementations of
// SecretReader and EnvWriter.
package sync
