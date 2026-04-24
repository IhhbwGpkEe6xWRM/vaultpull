// Package envscope provides prefix-based scoping for environment variable maps.
//
// A Scope declares which key prefixes are permitted. Keys that do not match
// any declared prefix are considered out-of-scope and can be filtered or
// blocked at write time.
//
// # Basic usage
//
//	s, _ := envscope.New([]string{"APP", "DB"})
//
//	// Filter a map down to in-scope keys only.
//	safe := s.Filter(raw)
//
//	// Validate that every key in a map is in-scope.
//	if err := s.Validate(raw); err != nil {
//	    log.Fatal(err)
//	}
//
// # Guard
//
// Guard wraps a Scope and intercepts individual or batch writes, recording
// which keys were blocked so callers can surface detailed diagnostics.
package envscope
