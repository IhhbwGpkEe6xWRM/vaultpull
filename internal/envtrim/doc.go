// Package envtrim provides a configurable Trimmer that removes unwanted
// entries from a secret map before it is written to a .env file.
//
// Three built-in rules are available:
//
//   - WithTrimEmpty    – removes keys whose value is the empty string.
//   - WithTrimWhitespace – removes keys whose value is blank (spaces/tabs).
//   - WithTrimFunc     – removes keys matching a caller-supplied predicate.
//
// Rules are additive; a value is removed if any rule matches.
//
// Example:
//
//	trimmer := envtrim.New(
//		envtrim.WithTrimEmpty(),
//		envtrim.WithTrimWhitespace(),
//	)
//	clean := trimmer.Apply(secrets)
package envtrim
