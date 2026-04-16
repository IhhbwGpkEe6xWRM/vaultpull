// Package expire implements TTL-based expiry checking for secrets fetched
// from Vault. It allows callers to determine whether a cached secret is
// still fresh, approaching expiry, or fully expired according to a
// configurable Policy.
//
// Usage:
//
//	p := expire.DefaultPolicy()
//	c := expire.New(p)
//	status := c.Check(lastFetchedAt)
//	if status == expire.StatusExpired {
//		// re-fetch from Vault
//	}
package expire
