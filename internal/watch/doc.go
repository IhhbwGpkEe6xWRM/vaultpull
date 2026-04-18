// Package watch provides secret polling and change detection for vaultpull.
//
// Use New to create a Watcher that periodically reads a Vault secret path and
// emits an Event on a channel whenever the returned secrets differ from the
// previously observed snapshot.
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
package watch
