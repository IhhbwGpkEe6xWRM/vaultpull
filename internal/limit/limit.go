// Package limit provides a rate limiter for Vault API requests,
// preventing thundering-herd issues when syncing many secrets.
package limit

import (
	"context"
	"time"
)

// Limiter controls the rate of outgoing requests.
type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	stop   chan struct{}
}

// Config holds rate-limiter settings.
type Config struct {
	// RequestsPerSecond is the maximum number of requests allowed per second.
	RequestsPerSecond int
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{RequestsPerSecond: 10}
}

// New creates a Limiter from the given Config and starts its background
// token-refill goroutine. Call Stop when done.
func New(cfg Config) *Limiter {
	rps := cfg.RequestsPerSecond
	if rps <= 0 {
		rps = 1
	}
	interval := time.Second / time.Duration(rps)
	l := &Limiter{
		ticker: time.NewTicker(interval),
		tokens: make(chan struct{}, rps),
		stop:   make(chan struct{}),
	}
	// Pre-fill one token so the first call is never blocked.
	l.tokens <- struct{}{}
	go l.refill()
	return l
}

func (l *Limiter) refill() {
	for {
		select {
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket full — drop the token
			}
		case <-l.stop:
			return
		}
	}
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop shuts down the background goroutine.
func (l *Limiter) Stop() {
	l.ticker.Stop()
	close(l.stop)
}
