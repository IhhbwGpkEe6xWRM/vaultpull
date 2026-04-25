// Package envwhistle scans environment variable maps for keys that match
// well-known sensitive patterns and reports structured findings.
//
// It is intended as a lightweight, non-blocking advisory layer: callers
// decide whether to fail, warn, or silently log based on the findings.
//
// Basic usage:
//
//	detector := envwhistle.New()
//	findings := detector.Scan(secrets)
//	if envwhistle.HasHigh(findings) {
//		// abort or prompt the user
//	}
//
//	reporter := envwhistle.NewReporter(os.Stderr)
//	reporter.Write(findings)
//	fmt.Println(envwhistle.Summary(findings))
//
// Custom rules can be registered via NewWithRules for project-specific
// naming conventions that should trigger warnings.
package envwhistle
