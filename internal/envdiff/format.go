package envdiff

import (
	"fmt"
	"io"
	"strings"
)

const (
	symAdd    = "+"
	symRemove = "-"
	symModify = "~"
	symSame   = "="
)

// Format writes a human-readable summary of r to w.
// When colour is true, ANSI escape codes are added.
func Format(w io.Writer, r Result, colour bool) {
	print := func(sym, key string) {
		if colour {
			var code string
			switch sym {
			case symAdd:
				code = "\033[32m"
			case symRemove:
				code = "\033[31m"
			case symModify:
				code = "\033[33m"
			default:
				code = "\033[0m"
			}
			fmt.Fprintf(w, "%s%s %s\033[0m\n", code, sym, key)
		} else {
			fmt.Fprintf(w, "%s %s\n", sym, key)
		}
	}
	for _, k := range r.Added {
		print(symAdd, k)
	}
	for _, k := range r.Removed {
		print(symRemove, k)
	}
	for _, k := range r.Modified {
		print(symModify, k)
	}
	if !r.HasChanges() {
		fmt.Fprintln(w, "no changes detected")
		return
	}
	parts := []string{}
	if n := len(r.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(r.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(r.Modified); n > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", n))
	}
	fmt.Fprintf(w, "summary: %s\n", strings.Join(parts, ", "))
}
