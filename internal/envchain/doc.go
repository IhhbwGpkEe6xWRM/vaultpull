// Package envchain provides a priority-ordered chain of secret readers.
//
// When multiple sources are available (e.g. Vault KV v1, Vault KV v2, a local
// .env file, or environment variables), envchain lets callers compose them
// into a single logical reader. Sources earlier in the chain take precedence:
// if the first reader supplies a non-empty value for a key, subsequent readers
// are not consulted for that key.
//
// Usage:
//
//	chain := envchain.New(primaryVault, fallbackVault, localFile)
//	secrets, err := chain.Resolve(ctx, "secret/myapp")
//
// Empty string values are treated as absent, allowing a downstream reader to
// supply a default without explicitly removing the key from an upstream map.
package envchain
