// Package validate provides pre-write validation for secret maps produced
// by the Vault client before they are persisted to .env files.
//
// Usage:
//
//	result := validate.Secrets(secrets)
//	if !result.OK() {
//	    log.Fatal(result.Err())
//	}
//
// Validation rules:
//   - Keys must not be empty.
//   - Keys must consist solely of ASCII letters, digits, and underscores.
//   - Values must not exceed 65 536 bytes.
package validate
