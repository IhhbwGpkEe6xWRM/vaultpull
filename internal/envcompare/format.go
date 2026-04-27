package envcompare

import (
	"fmt"
	"io"
	"strings"
)

const (
	symMatch    = "="
	symMismatch = "≠"
	symLeftOnly = "<"
	symRightOnly = ">"
)

// Format writes a human-readable comparison table to w.
// leftLabel and rightLabel are column headers (e.g. "vault", "local").
func Format(w io.Writer, r Result, leftLabel, rightLabel string) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "no keys to compare")
		return
	}

	keyW := len("KEY")
	for _, e := range r.Entries {
		if len(e.Key) > keyW {
			keyW = len(e.Key)
		}
	}

	header := fmt.Sprintf("%-*s  SYM  %-20s  %-20s", keyW, "KEY", leftLabel, rightLabel)
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, e := range r.Entries {
		sym := symbolFor(e.Status)
		lv := truncate(e.Left, 20)
		rv := truncate(e.Right, 20)
		fmt.Fprintf(w, "%-*s  %-3s  %-20s  %-20s\n", keyW, e.Key, sym, lv, rv)
	}

	counts := summaryCounts(r)
	fmt.Fprintf(w, "\n%d match, %d mismatch, %d left-only, %d right-only\n",
		counts[StatusMatch], counts[StatusMismatch],
		counts[StatusLeftOnly], counts[StatusRightOnly])
}

func symbolFor(s Status) string {
	switch s {
	case StatusMatch:
		return symMatch
	case StatusMismatch:
		return symMismatch
	case StatusLeftOnly:
		return symLeftOnly
	case StatusRightOnly:
		return symRightOnly
	default:
		return "?"
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func summaryCounts(r Result) map[Status]int {
	m := map[Status]int{
		StatusMatch:     0,
		StatusMismatch:  0,
		StatusLeftOnly:  0,
		StatusRightOnly: 0,
	}
	for _, e := range r.Entries {
		m[e.Status]++
	}
	return m
}
