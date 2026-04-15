// Package output provides a lightweight Formatter for writing structured
// status messages and sync summaries to stdout and stderr.
//
// It supports a quiet mode that suppresses informational output while
// still surfacing warnings and errors — suitable for use in scripts or
// CI pipelines where only failures should produce output.
//
// Usage:
//
//	f := output.New(quiet)
//	f.Info("reading secrets from vault...")
//	f.Summary(".env", len(secrets), cfg.Namespace != "")
//
// For testing, use NewWithWriters to inject custom io.Writer instances.
package output
