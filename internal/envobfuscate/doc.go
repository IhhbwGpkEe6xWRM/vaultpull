// Package envobfuscate replaces environment variable key names with
// deterministic HMAC-SHA256 aliases derived from a caller-supplied salt.
//
// This is useful when secret key names themselves are sensitive (e.g. they
// reveal the name of an internal service) and must not appear in logs,
// audit trails, or intermediate storage.
//
// Usage:
//
//	o := envobfuscate.New([]byte("my-salt"))
//	obfuscated := o.Apply(secrets)   // key names replaced with hex aliases
//	reveal := o.Reveal(secrets)      // alias -> original key lookup table
//
// The alias length is fixed at 16 hex characters (64-bit prefix of the
// HMAC digest). The same salt and key always produce the same alias, so
// the transformation is fully deterministic and reversible by any party
// that holds the salt.
package envobfuscate
