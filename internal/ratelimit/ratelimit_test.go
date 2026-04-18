package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ratelimit"
)

func TestWait_FirstCall_NoDelay(t *testing.T) {
	th := ratelimit.New(500 * time.Millisecond)
	start := time.Now()
	if err := th.Wait(context.Background(), "secret/foo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("first call should not block, took %v", elapsed)
	}
}

func TestWait_SecondCall_Throttled(t *testing.T) {
	th := ratelimit.New(100 * time.Millisecond)
	_ = th.Wait(context.Background(), "secret/bar")
	start := time.Now()
	_ = th.Wait(context.Background(), "secret/bar")
	if elapsed := time.Since(start); elapsed < 80*time.Millisecond {
		t.Errorf("second call should be throttled, elapsed %v", elapsed)
	}
}

func TestWait_DifferentPaths_Independent(t *testing.T) {
	th := ratelimit.New(500 * time.Millisecond)
	_ = th.Wait(context.Background(), "secret/a")
	start := time.Now()
	if err := th.Wait(context.Background(), "secret/b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("different path should not be throttled, took %v", elapsed)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	th := ratelimit.New(2 * time.Second)
	_ = th.Wait(context.Background(), "secret/x")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := th.Wait(ctx, "secret/x")
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestReset_ClearsThrottle(t *testing.T) {
	th := ratelimit.New(500 * time.Millisecond)
	_ = th.Wait(context.Background(), "secret/reset")
	th.Reset("secret/reset")
	start := time.Now()
	_ = th.Wait(context.Background(), "secret/reset")
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("after reset, call should not block, took %v", elapsed)
	}
}

func TestNew_ZeroInterval_UsesDefault(t *testing.T) {
	th := ratelimit.New(0)
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}
