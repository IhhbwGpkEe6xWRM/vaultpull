// Package envpromote provides utilities for promoting environment-specific
// secret sets from one environment to another (e.g. staging → production),
// with optional key filtering and dry-run support.
package envpromote

import "fmt"

// Option configures a Promoter.
type Option func(*Promoter)

// WithDryRun returns an Option that enables dry-run mode; no writes occur.
func WithDryRun() Option {
	return func(p *Promoter) { p.dryRun = true }
}

// WithAllowList returns an Option that restricts promotion to the given keys.
func WithAllowList(keys []string) Option {
	return func(p *Promoter) {
		p.allowList = make(map[string]struct{}, len(keys))
		for _, k := range keys {
			p.allowList[k] = struct{}{}
		}
	}
}

// Result describes the outcome of a promotion operation.
type Result struct {
	Promoted []string
	Skipped  []string
	DryRun   bool
}

// Promoter copies secrets from a source map into a destination map.
type Promoter struct {
	allowList map[string]struct{}
	dryRun    bool
}

// New creates a Promoter with the supplied options.
func New(opts ...Option) *Promoter {
	p := &Promoter{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Promote merges src into dst, respecting the allow-list and dry-run flag.
// dst is never mutated when dry-run is enabled.
func (p *Promoter) Promote(src, dst map[string]string) (Result, error) {
	if src == nil {
		return Result{}, fmt.Errorf("envpromote: source map must not be nil")
	}
	if dst == nil && !p.dryRun {
		return Result{}, fmt.Errorf("envpromote: destination map must not be nil")
	}

	var res Result
	res.DryRun = p.dryRun

	for k, v := range src {
		if len(p.allowList) > 0 {
			if _, ok := p.allowList[k]; !ok {
				res.Skipped = append(res.Skipped, k)
				continue
			}
		}
		if !p.dryRun {
			dst[k] = v
		}
		res.Promoted = append(res.Promoted, k)
	}

	sortStrings(res.Promoted)
	sortStrings(res.Skipped)
	return res, nil
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
