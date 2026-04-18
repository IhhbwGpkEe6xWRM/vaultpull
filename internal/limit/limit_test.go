package limit

import (
	"context"
	"testing"
	"time"
)

func TestWait_AcquiresToken(t *testing.T) {
	l := New(Config{RequestsPerSecond: 10})
	defer l.Stop()

	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	// Use RPS=1 and drain the pre-filled token so the next Wait must block.
	l := New(Config{RequestsPerSecond: 1})
	defer l.Stop()

	// Drain the pre-filled token.
	_ = l.Wait(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.RequestsPerSecond != 10 {
		t.Fatalf("expected 10 rps, got %d", cfg.RequestsPerSecond)
	}
}

func TestNew_ZeroRPS_DefaultsToOne(t *testing.T) {
	l := New(Config{RequestsPerSecond: 0})
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should still be able to acquire the pre-filled token.
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_MultipleTokens(t *testing.T) {
	l := New(Config{RequestsPerSecond: 50})
	defer l.Stop()

	// Give the refill goroutine time to add more tokens.
	time.Sleep(60 * time.Millisecond)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
}
