// Package watch provides secret polling and change detection for vaultpull.
//
// Use New to create a Watcher that periodically reads a Vault secret path and
// emits an Event on a channel whenever the returned secrets differ from the
// previously observed snapshot.
//
// # Event fields
//
// Each Event carries the secret Path, the new Data map, and a Timestamp
// indicating when the change was detected.
//
// # Debouncing
//
// For high-frequency paths, wrap the event channel with a Debouncer to
// coalesce bursts of changes into a single notification after a quiet period:
//
//	raw := make(chan watch.Event, 8)
//	debounced := make(chan watch.Event, 8)
//
//	w := watch.New(cfg, vaultClient)
//	go w.Watch(ctx, raw)
//
//	d := watch.NewDebouncer(2 * time.Second)
//	go d.Run(raw, debounced)
//
//	for ev := range debounced {
//		// handle ev
//	}
//
// # Context cancellation
//
// Both Watcher.Watch and Debouncer.Run respect context cancellation: when the
// provided context is cancelled, they stop processing and close their output
// channels cleanly.
package watch
