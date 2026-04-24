// Package envsplit partitions a flat secret map into named groups using
// prefix-based rules.
//
// A common pattern when pulling secrets from Vault is to store all secrets
// for a service under a single path but with prefixed keys, e.g.:
//
//	DB_HOST=localhost
//	DB_PORT=5432
//	API_KEY=abc123
//
// envsplit lets callers define rules that split such a map into logical
// groups, each with the prefix stripped:
//
//	s, _ := envsplit.New([]envsplit.Rule{
//	    {Name: "db",  Prefix: "DB_"},
//	    {Name: "api", Prefix: "API_"},
//	})
//	res := s.Split(secrets)
//	// res.Groups["db"]  => {"HOST": "localhost", "PORT": "5432"}
//	// res.Groups["api"] => {"KEY": "abc123"}
//
// Keys that do not match any rule are collected in Result.Remainder.
// Prefix matching is case-insensitive; first matching rule wins.
package envsplit
