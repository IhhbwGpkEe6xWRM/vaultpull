// Package envexpand resolves shell-style variable references within secret
// values before they are written to a .env file.
//
// References use the standard ${KEY} or $KEY syntax. Resolution is performed
// against the secret map itself, with an optional fallback to the host
// process environment via WithOSFallback.
//
// Self-referential expansions (e.g. FOO=${FOO}) are detected and the
// reference is replaced with an empty string to avoid infinite loops.
//
// Usage:
//
//	e := envexpand.New(envexpand.WithOSFallback())
//	resolved := e.Expand(secrets)
package envexpand
