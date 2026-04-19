// Package envencrypt provides AES-GCM encryption and decryption
// for Vault secret values, allowing secrets to be stored encrypted
// in .env files and decrypted at runtime.
//
// Usage:
//
//	key, err := envencrypt.LoadKeyFile("/etc/vaultpull/enc.key")
//	if err != nil { ... }
//
//	enc, err := envencrypt.New(key)
//	if err != nil { ... }
//
//	ciphertexts, err := enc.EncryptMap(secrets)
//	if err != nil { ... }
//
//	plaintexts, err := enc.DecryptMap(ciphertexts)
//	if err != nil { ... }
//
// Keys may be loaded from a file (raw bytes or hex-encoded) or from
// an environment variable via KeyFromEnv.
package envencrypt
