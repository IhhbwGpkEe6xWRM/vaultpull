// Package vault provides a thin wrapper around the HashiCorp Vault API client
// tailored for vaultpull's use-case: reading KV secrets from a configured path
// with optional namespace prefixing.
//
// # Usage
//
//	client, err := vault.NewClient(address, token, namespace)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	secrets, err := client.ReadSecrets("secret/myapp")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The client transparently handles both KV v1 and KV v2 secret engines by
// detecting whether the response payload contains a nested "data" key.
//
// Namespace handling:
//
// When a namespace is provided (e.g. "team-a"), it is prepended to every
// secret path so that "secret/myapp" becomes "team-a/secret/myapp". Leading
// and trailing slashes on both the namespace and path are normalised
// automatically.
package vault
