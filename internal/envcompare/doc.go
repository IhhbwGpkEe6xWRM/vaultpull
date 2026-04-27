// Package envcompare performs side-by-side comparison of two secret maps.
//
// It is typically used to surface drift between a Vault secret path and a
// local .env file before overwriting, giving operators visibility into what
// will change.
//
// Basic usage:
//
//	c := envcompare.New(true) // mask=true hides actual secret values
//	result := c.Compare(vaultSecrets, localSecrets)
//	if !result.Matches() {
//		envcompare.Format(os.Stdout, result, "vault", "local")
//	}
//
// Each Entry in the result carries a Status:
//
//	- StatusMatch     both sides have the same value
//	- StatusMismatch  both sides have the key but different values
//	- StatusLeftOnly  key is present only in the left (vault) map
//	- StatusRightOnly key is present only in the right (local) map
package envcompare
