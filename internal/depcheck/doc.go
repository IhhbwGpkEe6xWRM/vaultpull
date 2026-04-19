// Package depcheck enforces a list of required secret keys before vaultpull
// writes output, ensuring that downstream applications always receive a
// complete configuration.
//
// Usage:
//
//	checker := depcheck.New([]string{"DB_HOST", "DB_PASS", "API_KEY"})
//	if violations := checker.Check(secrets); len(violations) > 0 {
//		fmt.Print(depcheck.Summary(violations))
//		os.Exit(1)
//	}
//
// Any key that is absent from the secrets map, or present with an empty
// string value, is reported as a Violation. Summary formats the full list
// into a human-readable message suitable for CLI output.
package depcheck
