// Package ratelimit provides per-path request throttling for Vault API calls.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Throttle tracks per-path rate limits using a simple token bucket.
type Throttle struct {
	mu       sync.Mutex
	buckets  map[string]time.Time
	interval time.Duration
}

// New returns a Throttle that enforces a minimum interval between calls per path.
func New(interval time.Duration) *Throttle {
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	return &Throttle{
		buckets:  make(map[string]time.Time),
		interval: interval,
	}
}

// Wait blocks until the path is ready to be called, or the context is cancelled.
func (t *Throttle) Wait(ctx context.Context, path string) error {
	for {
		delay := t.nextDelay(path)
		if delay <= 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func (t *Throttle) nextDelay(path string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	next, ok := t.buckets[path]
	if !ok || now.After(next) {
		t.buckets[path] = now.Add(t.interval)
		return 0
	}
	return next.Sub(now)
}

// Reset clears the throttle state for a given path.
func (t *Throttle) Reset(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, path)
}
