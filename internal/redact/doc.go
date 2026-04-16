// Package redact provides key-pattern-based redaction of secret maps.
//
// It is used to sanitise output before secrets are printed to the terminal
// or written to audit logs, ensuring that values associated with sensitive
// keys (passwords, tokens, private keys, etc.) are never exposed in
// plain text outside of the .env file itself.
//
// Usage:
//
//	r := redact.New()
//	safe := r.Redact(secrets) // values for sensitive keys become "[REDACTED]"
//
// Custom patterns and placeholders can be supplied via NewWithPatterns.
package redact
