// Package envmask provides selective masking of environment variable values
// based on configurable key patterns.
//
// By default, keys matching common sensitive patterns (password, secret, token,
// api_key, credential, private) have their values replaced with a placeholder
// string before the map is used for display or logging purposes.
//
// # Basic usage
//
//	m, _ := envmask.New()
//	safe := m.Apply(secrets)   // sensitive values replaced with "***"
//	keys := m.MaskedKeys(secrets) // list of keys that were masked
//
// # Custom patterns
//
//	m, _ := envmask.New(envmask.WithPatterns([]string{"(?i)internal"}))
//
// # Reporting
//
//	r := envmask.NewReporter(os.Stderr)
//	r.Write(m.MaskedKeys(secrets))
package envmask
