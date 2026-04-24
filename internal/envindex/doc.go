// Package envindex provides a reverse-lookup index over a secrets map.
//
// It allows callers to quickly determine which keys share the same value,
// making it straightforward to detect accidental secret reuse before writing
// an env file.
//
// Basic usage:
//
//	idx := envindex.New(secrets)
//	if idx.HasDuplicates() {
//		for val, keys := range idx.Duplicates() {
//			fmt.Printf("value shared by %v: %q\n", keys, val)
//		}
//	}
package envindex
