// Package tokenize provides path tokenization utilities for splitting,
// inspecting, and reconstructing Vault secret paths.
//
// A Tokenizer breaks a raw path string into ordered segments (Parts)
// using a configurable separator (default "/"). It supports depth
// calculation, parent path derivation, and path reconstruction via Join.
//
// Example:
//
//	tz := tokenize.New()
//	tok := tz.Parse("secret/myapp/database")
//	fmt.Println(tok.Parts)       // [secret myapp database]
//	fmt.Println(tz.Depth(tok))   // 3
//	parent, _ := tz.Parent(tok)
//	fmt.Println(tz.Join(parent)) // secret/myapp
package tokenize
