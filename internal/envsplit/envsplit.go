// Package envsplit splits a flat secret map into named groups based on
// key prefix rules. Each group collects keys whose names match a given
// prefix, stripping that prefix from the resulting keys.
package envsplit

import (
	"fmt"
	"strings"
)

// Rule defines a named group and the key prefix that populates it.
type Rule struct {
	// Name is the logical group identifier (e.g. "db", "api").
	Name string
	// Prefix is the key prefix used to select secrets into this group.
	// Matching is case-insensitive; the prefix is stripped from result keys.
	Prefix string
}

// Splitter partitions a secret map into labelled groups.
type Splitter struct {
	rules []Rule
}

// Result holds the output of a split operation.
type Result struct {
	// Groups maps each rule name to its filtered, prefix-stripped secrets.
	Groups map[string]map[string]string
	// Remainder contains keys that did not match any rule.
	Remainder map[string]string
}

// New creates a Splitter from the provided rules.
// Returns an error if any rule has an empty Name or Prefix.
func New(rules []Rule) (*Splitter, error) {
	for i, r := range rules {
		if strings.TrimSpace(r.Name) == "" {
			return nil, fmt.Errorf("rule %d: name must not be empty", i)
		}
		if strings.TrimSpace(r.Prefix) == "" {
			return nil, fmt.Errorf("rule %d (%q): prefix must not be empty", i, r.Name)
		}
	}
	return &Splitter{rules: rules}, nil
}

// Split partitions secrets according to the configured rules.
// A key may match at most one rule (first match wins).
func (s *Splitter) Split(secrets map[string]string) Result {
	groups := make(map[string]map[string]string, len(s.rules))
	for _, r := range s.rules {
		groups[r.Name] = make(map[string]string)
	}
	remainder := make(map[string]string)

outer:
	for k, v := range secrets {
		upper := strings.ToUpper(k)
		for _, r := range s.rules {
			pfx := strings.ToUpper(r.Prefix)
			if strings.HasPrefix(upper, pfx) {
				stripped := k[len(r.Prefix):]
				stripped = strings.TrimPrefix(stripped, "_")
				if stripped == "" {
					stripped = k
				}
				groups[r.Name][stripped] = v
				continue outer
			}
		}
		remainder[k] = v
	}
	return Result{Groups: groups, Remainder: remainder}
}
