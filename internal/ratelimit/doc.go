// Package ratelimit provides per-path request throttling to avoid overwhelming
// the Vault API during bulk secret reads.
//
// Usage:
//
//	throttle := ratelimit.New(200 * time.Millisecond)
//	if err := throttle.Wait(ctx, secretPath); err != nil {
//		return err
//	}
//	// safe to call Vault now
package ratelimit
