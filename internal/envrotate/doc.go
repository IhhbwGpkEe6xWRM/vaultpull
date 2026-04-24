// Package envrotate applies a caller-supplied rotation function to every
// secret in a map, returning the updated secrets together with a per-key
// Result that records whether the value actually changed.
//
// Typical usage:
//
//	fn := func(key, old string) (string, error) {
//		// fetch fresh value from Vault or generate a new credential
//		return vaultClient.Read(key)
//	}
//
//	rotator, err := envrotate.New(fn)
//	if err != nil { ... }
//
//	updated, results, err := rotator.Apply(currentSecrets)
//	for _, r := range results {
//		if r.Rotated {
//			log.Printf("rotated %s (%s → %s)", r.Key, r.OldHash, r.NewHash)
//		}
//	}
//
// Values are never written to logs; only an 8-character djb2 hash is exposed
// so that callers can detect changes without leaking sensitive data.
package envrotate
