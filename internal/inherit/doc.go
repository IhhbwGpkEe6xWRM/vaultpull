// Package inherit implements secret inheritance for Vault paths.
//
// When secrets are organised in a hierarchy such as:
//
//	app/           → base secrets shared by all services
//	app/backend/   → secrets for the backend service
//	app/backend/prod → production overrides
//
// Resolver.Resolve walks from the shallowest ancestor down to the
// requested path, merging each level's secrets. Values from deeper
// paths take precedence over shallower ones, allowing teams to define
// defaults at a high level while overriding specific keys closer to
// the leaf.
//
// Usage:
//
//	r := inherit.New(vaultClient.ReadSecrets)
//	secrets, err := r.Resolve("app/backend/prod")
package inherit
