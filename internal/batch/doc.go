// Package batch provides concurrent fetching of secrets from multiple
// Vault paths. It is useful when a sync operation targets several paths
// and sequential reads would be too slow.
//
// Basic usage:
//
//	f := batch.New(vaultClient, 8)
//	results := f.FetchAll(ctx, []string{"secret/app", "secret/db"})
//	merged, err := batch.Merge(results)
//
// The concurrency parameter controls how many Vault reads run in
// parallel. A value of zero or below defaults to 4.
//
// Merge combines all successful results into one map. If any result
// carries an error, Merge returns that error immediately.
package batch
