// Package resolve translates logical secret paths into full Vault HTTP API
// paths, accounting for KV secrets engine version differences.
//
// KV v1 paths follow the form:
//
//	<mount>/<secret-path>
//
// KV v2 paths follow the form:
//
//	<mount>/data/<secret-path>   (read/write)
//	<mount>/metadata/<secret-path> (list/metadata)
//
// Example:
//
//	r, err := resolve.New("secret", resolve.KVv2)
//	if err != nil { ... }
//	apiPath := r.DataPath("myapp/production") // "secret/data/myapp/production"
package resolve
