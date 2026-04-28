// Package envaudit provides structured audit logging for environment secret
// operations. It records actions such as reads, writes, and deletions of
// individual secret keys, along with timestamps, paths, and optional metadata.
//
// Each audit entry is written as a newline-delimited JSON record, making the
// output suitable for ingestion by log aggregators or post-processing tools.
//
// Basic usage:
//
//	auditor, err := envaudit.New("/var/log/vaultpull/audit.jsonl")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer auditor.Close()
//
//	auditor.Record(ctx, envaudit.ActionRead, "secret/app", "DB_PASSWORD")
//
// The WithClock option allows injecting a custom clock for deterministic
// testing without relying on real wall-clock time.
package envaudit
