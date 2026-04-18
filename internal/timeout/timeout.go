// Package timeout wraps a context with a configurable deadline and
// provides helpers for common Vault operation timeouts.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// DefaultRequestTimeout is used when no explicit duration is configured.
const DefaultRequestTimeout = 10 * time.Second

// ErrDeadlineExceeded is returned when an operation exceeds its deadline.
var ErrDeadlineExceeded = errors.New("operation timed out")

// Config holds timeout settings.
type Config struct {
	// Request is the per-request deadline. Zero means DefaultRequestTimeout.
	Request time.Duration
}

// Limiter wraps operations with a deadline derived from Config.
type Limiter struct {
	cfg Config
}

// New returns a Limiter using the provided Config.
// If cfg.Request is zero it is set to DefaultRequestTimeout.
func New(cfg Config) *Limiter {
	if cfg.Request <= 0 {
		cfg.Request = DefaultRequestTimeout
	}
	return &Limiter{cfg: cfg}
}

// Do runs fn within a deadline context derived from parent.
// If fn does not return before the deadline, ErrDeadlineExceeded is returned.
func (l *Limiter) Do(parent context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(parent, l.cfg.Request)
	defer cancel()

	err := fn(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("%w: %w", ErrDeadlineExceeded, err)
	}
	return err
}
