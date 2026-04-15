// Package retry provides configurable retry logic with exponential backoff
// for transient errors encountered during Vault API calls.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Config holds the retry policy configuration.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential backoff delay.
	MaxDelay time.Duration
	// Multiplier is the factor applied to the delay on each attempt.
	Multiplier float64
}

// DefaultConfig returns a sensible default retry policy.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  4,
		InitialDelay: 250 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Do executes fn up to cfg.MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns a non-retryable error.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		var nr *NonRetryableError
		if errors.As(lastErr, &nr) {
			return lastErr
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay = time.Duration(math.Min(
			float64(delay)*cfg.Multiplier,
			float64(cfg.MaxDelay),
		))
	}

	return errors.Join(ErrMaxAttemptsReached, lastErr)
}

// NonRetryableError wraps an error that should not trigger a retry.
type NonRetryableError struct {
	Cause error
}

func (e *NonRetryableError) Error() string {
	return "non-retryable: " + e.Cause.Error()
}

func (e *NonRetryableError) Unwrap() error { return e.Cause }

// Permanent wraps err so that Do stops immediately without retrying.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &NonRetryableError{Cause: err}
}
