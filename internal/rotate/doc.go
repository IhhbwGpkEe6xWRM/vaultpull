// Package rotate detects secret rotation by comparing a current set of
// Vault secrets against a previously recorded snapshot.
//
// Usage:
//
//	detector := rotate.NewDetector("secret/myapp")
//	events := detector.Detect(previousSecrets, currentSecrets)
//	for _, e := range events {
//		fmt.Printf("rotated: %s at %s\n", e.Key, e.DetectedAt)
//	}
//
// A rotation event is recorded only when a key existed in the previous
// snapshot and its value differs in the current fetch. Newly added keys
// are not considered rotations.
package rotate
