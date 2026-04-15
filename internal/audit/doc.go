// Package audit provides append-only structured audit logging for vaultpull.
//
// Each sync operation records audit entries describing which Vault secret
// paths were synced, skipped due to namespace filtering, or failed. Entries
// are written as newline-delimited JSON to a configurable file path.
//
// Usage:
//
//	logger, err := audit.NewLogger("/var/log/vaultpull-audit.log")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	_ = logger.Record(audit.EventSynced, "secret/myapp/db", "")
//	_ = logger.Record(audit.EventSkipped, "secret/other/key", "namespace mismatch")
//
// Pass an empty path to NewLogger to disable file logging entirely;
// the returned logger discards all writes without error.
package audit
