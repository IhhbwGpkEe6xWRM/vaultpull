package rotate

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Notifier writes rotation event summaries to a writer.
type Notifier struct {
	w io.Writer
}

// NewNotifier creates a Notifier that writes to stderr by default.
func NewNotifier() *Notifier {
	return &Notifier{w: os.Stderr}
}

// NewNotifierWithWriter creates a Notifier using the provided writer.
func NewNotifierWithWriter(w io.Writer) *Notifier {
	return &Notifier{w: w}
}

// Notify prints each rotation event to the writer. It returns the number
// of events reported and any write error encountered.
func (n *Notifier) Notify(events []Event) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[rotate] %s\n", Summary(events)))
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("  ~ %s (%s) at %s\n",
			e.Key, e.Path, e.DetectedAt.Format("2006-01-02T15:04:05Z")))
	}
	_, err := fmt.Fprint(n.w, sb.String())
	return len(events), err
}
