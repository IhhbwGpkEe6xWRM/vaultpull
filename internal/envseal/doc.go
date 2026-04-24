// Package envseal provides HMAC-SHA256 tamper detection for secret maps written
// to .env files by vaultpull.
//
// # Overview
//
// After writing secrets to a .env file, callers can seal the contents:
//
//	sealer, _ := envseal.New(key)
//	_ = sealer.WriteSealFile(secrets, ".env.seal")
//
// On the next sync, verify the file has not been modified outside of vaultpull:
//
//	if err := sealer.VerifySealFile(secrets, ".env.seal"); err != nil {
//		// ErrTampered or ErrMissingSeal
//	}
//
// # Key Loading
//
// Keys can be loaded from a file (raw bytes or hex-encoded) or from an
// environment variable using LoadKeyFile and KeyFromEnv respectively.
//
// # Signature Format
//
// The seal file contains a single hex-encoded HMAC-SHA256 digest computed over
// the lexicographically sorted "key=value" pairs of the secret map.
package envseal
