// Package retry implements a simple exponential-backoff retry helper used
// throughout vaultpull to handle transient Vault API failures gracefully.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return callVault()
//	})
//
// To prevent a specific error from being retried, wrap it with retry.Permanent:
//
//	return retry.Permanent(fmt.Errorf("permission denied"))
//
// The package distinguishes between:
//   - Transient errors: retried up to MaxAttempts with exponential backoff.
//   - Permanent errors: returned immediately without further attempts.
//   - Context cancellation: respected between attempts.
package retry
