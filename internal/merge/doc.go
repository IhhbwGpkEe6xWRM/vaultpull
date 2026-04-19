// Package merge combines secret maps from multiple Vault paths or sources
// into a single map using a configurable conflict resolution strategy.
//
// Three strategies are available:
//
//   - LastWins: the last source to define a key takes precedence.
//   - FirstWins: the first source to define a key takes precedence.
//   - ErrorOnConflict: returns an error if any key is defined more than once.
//
// All conflicts are recorded in the Result regardless of strategy (except
// ErrorOnConflict, which halts processing immediately).
//
// # Example
//
//	result, err := merge.Merge(merge.LastWins, maps...)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Merged)
package merge
