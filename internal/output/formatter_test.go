package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/output"
)

func newTestFormatter(quiet bool) (*output.Formatter, *bytes.Buffer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	f := output.NewWithWriters(out, err, quiet)
	return f, out, err
}

func TestInfo_WritesToStdout(t *testing.T) {
	f, out, _ := newTestFormatter(false)
	f.Info("loading secrets")
	if !strings.Contains(out.String(), "loading secrets") {
		t.Errorf("expected stdout to contain message, got: %q", out.String())
	}
}

func TestInfo_QuietMode_NoOutput(t *testing.T) {
	f, out, _ := newTestFormatter(true)
	f.Info("should be silent")
	if out.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %q", out.String())
	}
}

func TestSuccess_ContainsCheckmark(t *testing.T) {
	f, out, _ := newTestFormatter(false)
	f.Success("done")
	if !strings.Contains(out.String(), "✓") {
		t.Errorf("expected success marker in output, got: %q", out.String())
	}
}

func TestWarn_WritesToStderr(t *testing.T) {
	f, out, err := newTestFormatter(false)
	f.Warn("something odd")
	if err.Len() == 0 {
		t.Error("expected stderr to have content")
	}
	if out.Len() != 0 {
		t.Error("expected stdout to be empty for warnings")
	}
}

func TestError_WritesToStderr(t *testing.T) {
	f, _, err := newTestFormatter(false)
	f.Error("vault unreachable")
	if !strings.Contains(err.String(), "vault unreachable") {
		t.Errorf("expected error in stderr, got: %q", err.String())
	}
}

func TestSummary_IncludesKeyCount(t *testing.T) {
	f, out, _ := newTestFormatter(false)
	f.Summary(".env", 5, false)
	if !strings.Contains(out.String(), "5 key(s)") {
		t.Errorf("expected key count in summary, got: %q", out.String())
	}
}

func TestSummary_WithFilter_MentionsNamespace(t *testing.T) {
	f, out, _ := newTestFormatter(false)
	f.Summary(".env", 3, true)
	if !strings.Contains(out.String(), "namespace filter") {
		t.Errorf("expected namespace filter note in summary, got: %q", out.String())
	}
}

func TestSummary_QuietMode_NoOutput(t *testing.T) {
	f, out, _ := newTestFormatter(true)
	f.Summary(".env", 10, true)
	if out.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %q", out.String())
	}
}
