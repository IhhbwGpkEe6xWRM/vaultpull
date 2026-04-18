// Package environ provides utilities for loading environment variables
// into secret maps and merging them with secrets fetched from Vault.
//
// It supports optional prefix filtering so only a scoped subset of the
// host environment is considered, and a configurable override flag that
// controls whether local env values win over Vault-sourced ones.
package environ
