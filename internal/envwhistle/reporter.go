package envwhistle

import (
	"fmt"
	"io"
	"strings"
)

// Reporter formats and writes Findings to an io.Writer.
type Reporter struct {
	w      io.Writer
	prefix string
}

// NewReporter returns a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w, prefix: "[whistle]"}
}

// NewReporterWithPrefix returns a Reporter with a custom log prefix.
func NewReporterWithPrefix(w io.Writer, prefix string) *Reporter {
	return &Reporter{w: w, prefix: prefix}
}

// Write emits all findings to the writer, one per line.
// Returns the number of findings written and any write error.
func (r *Reporter) Write(findings []Finding) (int, error) {
	if len(findings) == 0 {
		return 0, nil
	}
	var errs []string
	for _, f := range findings {
		line := fmt.Sprintf("%s [%s] %s — %s\n", r.prefix, strings.ToUpper(string(f.Severity)), f.Key, f.Reason)
		if _, err := fmt.Fprint(r.w, line); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return len(findings) - len(errs), fmt.Errorf("write errors: %s", strings.Join(errs, "; "))
	}
	return len(findings), nil
}

// Summary returns a single-line summary of findings by severity.
func Summary(findings []Finding) string {
	counts := map[Severity]int{}
	for _, f := range findings {
		counts[f.Severity]++
	}
	if len(counts) == 0 {
		return "no findings"
	}
	parts := []string{}
	for _, sev := range []Severity{SeverityHigh, SeverityMedium, SeverityLow} {
		if n := counts[sev]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, sev))
		}
	}
	return strings.Join(parts, ", ")
}
