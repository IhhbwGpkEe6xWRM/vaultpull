// Package envlock implements lightweight file-based locking for .env output
// files written by vaultpull.
//
// When multiple vaultpull processes run concurrently against the same output
// path, envlock prevents interleaved writes by serialising access through a
// sibling lock file (.<filename>.lock).
//
// Usage:
//
//	locker := envlock.New(5 * time.Second)
//	lk, err := locker.Acquire("/path/to/.env")
//	if err != nil {
//		// another process holds the lock or timeout exceeded
//	}
//	defer lk.Release()
package envlock
