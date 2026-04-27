package envmask_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envmask"
)

func TestWrite_NoKeys_WritesNothing(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporter(&buf)
	r.Write(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestWrite_ContainsKeyNames(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporter(&buf)
	r.Write([]string{"DB_PASSWORD", "TOKEN"})
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("output missing DB_PASSWORD: %q", out)
	}
	if !strings.Contains(out, "TOKEN") {
		t.Errorf("output missing TOKEN: %q", out)
	}
}

func TestWrite_ContainsCount(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporter(&buf)
	r.Write([]string{"A", "B", "C"})
	if !strings.Contains(buf.String(), "3") {
		t.Errorf("expected count 3 in output: %q", buf.String())
	}
}

func TestWrite_CustomPrefix(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporterWithPrefix(&buf, "[custom]")
	r.Write([]string{"SECRET"})
	if !strings.HasPrefix(buf.String(), "[custom]") {
		t.Errorf("expected custom prefix: %q", buf.String())
	}
}

func TestSummary_Empty(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporter(&buf)
	if got := r.Summary(nil); got != "no keys masked" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithKeys(t *testing.T) {
	var buf strings.Builder
	r := envmask.NewReporter(&buf)
	s := r.Summary([]string{"TOKEN", "SECRET"})
	if !strings.Contains(s, "2") {
		t.Errorf("expected count in summary: %q", s)
	}
	if !strings.Contains(s, "TOKEN") {
		t.Errorf("expected TOKEN in summary: %q", s)
	}
}
