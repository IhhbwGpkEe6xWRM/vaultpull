// Package output provides formatting utilities for displaying sync results
// and status information to the user via stdout/stderr.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Level represents the severity of a log message.
type Level int

const (
	LevelInfo  Level = iota
	LevelWarn
	LevelError
	LevelSuccess
)

// Formatter writes formatted output to configurable writers.
type Formatter struct {
	out io.Writer
	err io.Writer
	quiet bool
}

// New returns a Formatter that writes to stdout and stderr.
func New(quiet bool) *Formatter {
	return &Formatter{
		out: os.Stdout,
		err: os.Stderr,
		quiet: quiet,
	}
}

// NewWithWriters returns a Formatter with custom writers (useful for testing).
func NewWithWriters(out, err io.Writer, quiet bool) *Formatter {
	return &Formatter{out: out, err: err, quiet: quiet}
}

// Info prints an informational message unless quiet mode is enabled.
func (f *Formatter) Info(msg string) {
	if !f.quiet {
		fmt.Fprintf(f.out, "  %s\n", msg)
	}
}

// Success prints a success message unless quiet mode is enabled.
func (f *Formatter) Success(msg string) {
	if !f.quiet {
		fmt.Fprintf(f.out, "✓ %s\n", msg)
	}
}

// Warn always prints a warning message to stderr.
func (f *Formatter) Warn(msg string) {
	fmt.Fprintf(f.err, "! %s\n", msg)
}

// Error always prints an error message to stderr.
func (f *Formatter) Error(msg string) {
	fmt.Fprintf(f.err, "✗ %s\n", msg)
}

// Summary prints a sync summary with counts of written keys.
func (f *Formatter) Summary(path string, keyCount int, filtered bool) {
	if f.quiet {
		return
	}
	parts := []string{fmt.Sprintf("%d key(s) written to %s", keyCount, path)}
	if filtered {
		parts = append(parts, "(namespace filter applied)")
	}
	fmt.Fprintf(f.out, "✓ %s\n", strings.Join(parts, " "))
}
