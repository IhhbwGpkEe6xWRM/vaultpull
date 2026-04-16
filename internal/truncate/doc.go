// Package truncate provides length-limiting utilities for secret values
// displayed in CLI output, logs, and audit records.
//
// Use [New] for default behaviour (80-character limit) or [NewWithLimit]
// to specify a custom threshold. Truncated values are suffixed with "..."
// to indicate that content has been omitted.
//
// Example:
//
//	tr := truncate.New()
//	safe := tr.Map(secrets) // returns a copy with long values trimmed
package truncate
