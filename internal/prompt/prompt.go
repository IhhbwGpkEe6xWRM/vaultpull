// Package prompt provides interactive confirmation prompts for destructive
// operations such as overwriting existing .env files with changed secrets.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirmer asks the user to confirm an action.
type Confirmer struct {
	in  io.Reader
	out io.Writer
}

// New returns a Confirmer that reads from stdin and writes to stdout.
func New() *Confirmer {
	return &Confirmer{in: os.Stdin, out: os.Stdout}
}

// NewWithReadWriter returns a Confirmer backed by the provided reader/writer.
func NewWithReadWriter(in io.Reader, out io.Writer) *Confirmer {
	return &Confirmer{in: in, out: out}
}

// Confirm prints prompt and returns true if the user types "y" or "yes".
// Any other input (including empty) returns false.
func (c *Confirmer) Confirm(prompt string) (bool, error) {
	fmt.Fprintf(c.out, "%s [y/N]: ", prompt)

	scanner := bufio.NewScanner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("prompt: read error: %w", err)
		}
		// EOF — treat as "no"
		return false, nil
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// MustConfirm is like Confirm but panics on read errors.
func (c *Confirmer) MustConfirm(prompt string) bool {
	ok, err := c.Confirm(prompt)
	if err != nil {
		panic(err)
	}
	return ok
}
