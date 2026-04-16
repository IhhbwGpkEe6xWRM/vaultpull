package expire

import (
	"testing"
	"time"
)

func checkerAt(now time.Time, p Policy) *Checker {
	c := New(p)
	c.now = func() time.Time { return now }
	return c
}

func TestCheck_Fresh(t *testing.T) {
	p := DefaultPolicy()
	base := time.Now()
	c := checkerAt(base.Add(1*time.Hour), p)
	if got := c.Check(base); got != StatusFresh {
		t.Fatalf("expected fresh, got %s", got)
	}
}

func TestCheck_Warning(t *testing.T) {
	p := DefaultPolicy() // TTL=24h, WarnAt=2h
	base := time.Now()
	// 23 hours after fetch => within 1h of expiry, inside WarnAt window
	c := checkerAt(base.Add(23*time.Hour), p)
	if got := c.Check(base); got != StatusWarning {
		t.Fatalf("expected warning, got %s", got)
	}
}

func TestCheck_Expired(t *testing.T) {
	p := DefaultPolicy()
	base := time.Now()
	c := checkerAt(base.Add(25*time.Hour), p)
	if got := c.Check(base); got != StatusExpired {
		t.Fatalf("expected expired, got %s", got)
	}
}

func TestCheck_ExactlyAtTTL(t *testing.T) {
	p := DefaultPolicy()
	base := time.Now()
	c := checkerAt(base.Add(p.TTL), p)
	if got := c.Check(base); got != StatusExpired {
		t.Fatalf("expected expired at exact TTL boundary, got %s", got)
	}
}

func TestTimeUntilExpiry_Positive(t *testing.T) {
	p := DefaultPolicy()
	base := time.Now()
	c := checkerAt(base.Add(1*time.Hour), p)
	rem := c.TimeUntilExpiry(base)
	if rem <= 0 {
		t.Fatalf("expected positive remainder, got %v", rem)
	}
}

func TestTimeUntilExpiry_Negative(t *testing.T) {
	p := DefaultPolicy()
	base := time.Now()
	c := checkerAt(base.Add(30*time.Hour), p)
	rem := c.TimeUntilExpiry(base)
	if rem >= 0 {
		t.Fatalf("expected negative remainder for expired entry, got %v", rem)
	}
}

func TestStatus_String(t *testing.T) {
	cases := map[Status]string{
		StatusFresh:   "fresh",
		StatusWarning: "warning",
		StatusExpired: "expired",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, got, want)
		}
	}
}
