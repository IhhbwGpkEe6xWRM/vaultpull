// Package transform provides a composable pipeline for transforming
// secret key/value pairs before they are written to .env files.
//
// # Overview
//
// A Pipeline is constructed from one or more Transformer functions and
// applied to a map[string]string of secrets. Each Transformer receives
// the key and value produced by the previous step, enabling flexible,
// ordered processing.
//
// # Built-in Transformers
//
//   - UppercaseKeys  – converts all keys to UPPER_CASE
//   - PrefixKeys     – prepends a static string to every key
//   - TrimSpace      – strips leading/trailing whitespace from keys and values
//   - ReplaceHyphens – replaces hyphens in keys with underscores
//   - DropNonPrintable – removes non-printable runes from values
//
// # Usage
//
//	p := transform.New(
//	    transform.TrimSpace(),
//	    transform.ReplaceHyphens(),
//	    transform.UppercaseKeys(),
//	)
//	cleaned := p.Apply(rawSecrets)
package transform
