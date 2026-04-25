// Package envwhistle detects secrets in environment maps that match
// well-known sensitive patterns and emits structured warnings.
package envwhistle

import (
	"regexp"
	"sort"
	"strings"
)

// Severity indicates how critical a finding is.
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
)

// Finding describes a single detected issue.
type Finding struct {
	Key      string
	Severity Severity
	Reason   string
}

type rule struct {
	pattern  *regexp.Regexp
	severity Severity
	reason   string
}

// Detector scans env maps for sensitive or suspicious keys.
type Detector struct {
	rules []rule
}

var defaultRules = []rule{
	{regexp.MustCompile(`(?i)(password|passwd|secret|token|api_?key)`), SeverityHigh, "key name suggests a credential"},
	{regexp.MustCompile(`(?i)(private_?key|pem|cert)`), SeverityHigh, "key name suggests a private key or certificate"},
	{regexp.MustCompile(`(?i)(access_?key|auth)`), SeverityMedium, "key name suggests an access credential"},
	{regexp.MustCompile(`(?i)(url|dsn|endpoint).*`), SeverityLow, "key name may contain a sensitive URL"},
}

// New returns a Detector using the built-in rule set.
func New() *Detector {
	return &Detector{rules: defaultRules}
}

// NewWithRules returns a Detector using only the provided patterns.
func NewWithRules(patterns []string, severity Severity, reason string) (*Detector, error) {
	var rules []rule
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule{pattern: re, severity: severity, reason: reason})
	}
	return &Detector{rules: rules}, nil
}

// Scan returns all findings for the given secret map.
func (d *Detector) Scan(secrets map[string]string) []Finding {
	var findings []Finding
	for key := range secrets {
		for _, r := range d.rules {
			if r.pattern.MatchString(strings.ToLower(key)) {
				findings = append(findings, Finding{
					Key:      key,
					Severity: r.severity,
					Reason:   r.reason,
				})
				break
			}
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Key < findings[j].Key
	})
	return findings
}

// HasHigh returns true if any finding has SeverityHigh.
func HasHigh(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityHigh {
			return true
		}
	}
	return false
}
