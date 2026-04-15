// Package diff compares two maps of secret key-value pairs and produces a
// structured Result describing which keys were added, removed, modified, or
// left unchanged between a previously cached snapshot and the freshly fetched
// secrets.
//
// Typical usage:
//
//	old := cache.Load()          // map[string]string from local cache
//	next := vault.ReadSecrets()  // map[string]string from Vault
//	result := diff.Compare(old, next)
//	if result.HasChanges() {
//		a, r, m := result.Summary()
//		fmt.Printf("+%d -%d ~%d\n", a, r, m)
//	}
package diff
