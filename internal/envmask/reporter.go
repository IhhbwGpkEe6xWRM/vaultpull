package envmask

import (
	"fmt"
	"io"
	"strings"
)

// Reporter writes a human-readable summary of masked keys to an io.Writer.
type Reporter struct {
	w      io.Writer
	prefix string
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w, prefix: "[envmask]"}
}

// NewReporterWithPrefix creates a Reporter with a custom log prefix.
func NewReporterWithPrefix(w io.Writer, prefix string) *Reporter {
	return &Reporter{w: w, prefix: prefix}
}

// Write prints each masked key to the underlying writer.
// It is a no-op when keys is empty.
func (r *Reporter) Write(keys []string) {
	if len(keys) == 0 {
		return
	}
	fmt.Fprintf(r.w, "%s masked %d key(s): %s\n",
		r.prefix, len(keys), strings.Join(keys, ", "))
}

// Summary returns a single-line description of the masked keys.
func (r *Reporter) Summary(keys []string) string {
	if len(keys) == 0 {
		return "no keys masked"
	}
	return fmt.Sprintf("%d key(s) masked: %s", len(keys), strings.Join(keys, ", "))
}
