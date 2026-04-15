package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

func fastConfig() retry.Config {
	return retry.Config{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnTransientError(t *testing.T) {
	var calls atomic.Int32
	transient := errors.New("transient")

	err := retry.Do(context.Background(), fastConfig(), func() error {
		if calls.Add(1) < 3 {
			return transient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	var calls atomic.Int32
	persistent := errors.New("always fails")

	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls.Add(1)
		return persistent
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestDo_PermanentErrorStopsImmediately(t *testing.T) {
	calls := 0
	perm := errors.New("forbidden")

	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		return retry.Permanent(perm)
	})
	if !errors.Is(err, perm) {
		t.Fatalf("expected underlying error %v, got %v", perm, err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, fastConfig(), func() error {
		calls++
		return errors.New("would retry")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls after pre-cancelled ctx, got %d", calls)
	}
}

func TestPermanent_NilPassthrough(t *testing.T) {
	if retry.Permanent(nil) != nil {
		t.Fatal("Permanent(nil) should return nil")
	}
}
