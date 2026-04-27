// Package envstamp attaches metadata stamps (build version, timestamp, hostname)
// to a secret map as additional synthetic keys.
package envstamp

import (
	"fmt"
	"os"
	"time"
)

// Option configures the Stamper.
type Option func(*Stamper)

// WithVersion adds a BUILD_VERSION key with the given value.
func WithVersion(v string) Option {
	return func(s *Stamper) { s.version = v }
}

// WithTimestamp adds a STAMP_TIMESTAMP key using the given clock function.
func WithTimestamp(fn func() time.Time) Option {
	return func(s *Stamper) { s.clock = fn }
}

// WithHostname adds a STAMP_HOSTNAME key with the current hostname.
func WithHostname() Option {
	return func(s *Stamper) { s.hostname = true }
}

// WithPrefix sets a prefix applied to every stamp key.
func WithPrefix(p string) Option {
	return func(s *Stamper) { s.prefix = p }
}

// Stamper injects metadata keys into a secret map.
type Stamper struct {
	version  string
	clock    func() time.Time
	hostname bool
	prefix   string
}

// New returns a Stamper configured with the supplied options.
func New(opts ...Option) *Stamper {
	s := &Stamper{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a copy of m with stamp keys injected. The original map is
// never mutated.
func (s *Stamper) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}

	key := func(name string) string {
		if s.prefix != "" {
			return fmt.Sprintf("%s_%s", s.prefix, name)
		}
		return name
	}

	if s.version != "" {
		out[key("BUILD_VERSION")] = s.version
	}
	if s.clock != nil {
		out[key("STAMP_TIMESTAMP")] = s.clock().UTC().Format(time.RFC3339)
	}
	if s.hostname {
		h, err := os.Hostname()
		if err != nil {
			h = "unknown"
		}
		out[key("STAMP_HOSTNAME")] = h
	}
	return out
}

// Keys returns the set of stamp key names that would be injected given the
// current configuration (without a prefix).
func (s *Stamper) Keys() []string {
	var keys []string
	if s.version != "" {
		keys = append(keys, "BUILD_VERSION")
	}
	if s.clock != nil {
		keys = append(keys, "STAMP_TIMESTAMP")
	}
	if s.hostname {
		keys = append(keys, "STAMP_HOSTNAME")
	}
	return keys
}
