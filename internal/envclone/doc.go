// Package envclone implements deep-copy semantics for secret maps.
//
// A Cloner produces an independent copy of a map[string]string, ensuring
// that mutations to the clone do not affect the original and vice-versa.
//
// Optional key filtering lets callers restrict which entries are carried
// over to the clone, and an optional value hook allows in-flight
// transformation (e.g. masking, encoding) of every value during the copy.
//
// Example:
//
//	cloner := envclone.New(
//		envclone.WithKeyFilter(func(k string) bool {
//			return !strings.HasPrefix(k, "INTERNAL_")
//		}),
//		envclone.WithValueHook(func(k, v string) (string, error) {
//			return masker.Mask(v), nil
//		}),
//	)
//
//	safe, err := cloner.Clone(secrets)
package envclone
