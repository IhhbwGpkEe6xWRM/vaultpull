// Package envrewrite implements key rewriting for secret maps.
//
// Rules can rename or drop keys before they are written to .env files,
// allowing vault secret paths to be mapped to application-specific
// environment variable names without modifying the vault data itself.
//
// Example usage:
//
//	rw, err := envrewrite.New([]envrewrite.Rule{
//		{From: "DB_PASSWORD", To: "DATABASE_PASSWORD"},
//		{From: "INTERNAL_TOKEN", To: ""},  // drop
//	})
//	out := rw.Apply(secrets)
package envrewrite
