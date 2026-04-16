// Package expire provides TTL-based expiry checking for cached secrets.
package expire

import (
	"time"
)

// Policy defines expiry rules for secrets.
type Policy struct {
	TTL      time.Duration
	WarnAt   time.Duration // warn when this close to expiry
}

// DefaultPolicy returns a sensible default expiry policy.
func DefaultPolicy() Policy {
	return Policy{
		TTL:    24 * time.Hour,
		WarnAt: 2 * time.Hour,
	}
}

// Status represents the expiry state of a secret entry.
type Status int

const (
	StatusFresh   Status = iota
	StatusWarning        // within WarnAt window
	StatusExpired
)

// String returns a human-readable label for the status.
func (s Status) String() string {
	switch s {
	case StatusFresh:
		return "fresh"
	case StatusWarning:
		return "warning"
	case StatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// Checker evaluates expiry status against a policy.
type Checker struct {
	policy Policy
	now    func() time.Time
}

// New returns a Checker using the given policy.
func New(p Policy) *Checker {
	return &Checker{policy: p, now: time.Now}
}

// Check returns the expiry Status for an entry last fetched at fetchedAt.
func (c *Checker) Check(fetchedAt time.Time) Status {
	now := c.now()
	age := now.Sub(fetchedAt)
	if age >= c.policy.TTL {
		return StatusExpired
	}
	if age >= c.policy.TTL-c.policy.WarnAt {
		return StatusWarning
	}
	return StatusFresh
}

// TimeUntilExpiry returns the duration remaining before the entry expires.
func (c *Checker) TimeUntilExpiry(fetchedAt time.Time) time.Duration {
	expireAt := fetchedAt.Add(c.policy.TTL)
	return expireAt.Sub(c.now())
}
