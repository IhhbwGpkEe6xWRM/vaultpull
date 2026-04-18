// Package scope implements prefix-based path scoping for Vault secret reads.
//
// A Scope is constructed with a list of allowed path prefixes. Any secret path
// that does not fall under one of those prefixes is rejected, enabling
// fine-grained control over which parts of the Vault hierarchy a vaultpull
// invocation may access.
//
// Usage:
//
//	s := scope.New([]string{"secrets/production", "secrets/shared"})
//	if s.Allows(path) {
//		// read secret
//	}
//
// When no prefixes are provided all paths are permitted, preserving backwards-
// compatible behaviour for configurations that do not declare a scope.
package scope
