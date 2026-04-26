// Package envsubst performs shell-style variable substitution on env map values.
//
// References of the form ${KEY} or $KEY within a value are replaced with the
// corresponding value from the same map (or an optional OS env fallback).
//
// Substitution is applied in multiple passes (default 5) to support chained
// references such as C -> B -> A. Circular references resolve to empty strings
// once the source key is exhausted.
//
// Example:
//
//	s := envsubst.New(envsubst.WithOSEnv(osEnv))
//	out, err := s.Apply(secrets)
package envsubst
