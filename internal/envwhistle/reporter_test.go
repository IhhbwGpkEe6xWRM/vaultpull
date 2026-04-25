package envwhistle_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envwhistle"
)

func TestWrite_NoFindings_WritesNothing(t *testing.T) {
	var buf bytes.Buffer
	r := envwhistle.NewReporter(&buf)
	n, err := r.Write(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWrite_ContainsSeverityAndKey(t *testing.T) {
	var buf bytes.Buffer
	r := envwhistle.NewReporter(&buf)
	findings := []envwhistle.Finding{
		{Key: "DB_PASSWORD", Severity: envwhistle.SeverityHigh, Reason: "credential"},
	}
	n, err := r.Write(findings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "HIGH") {
		t.Errorf("expected HIGH in output, got %q", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key in output, got %q", out)
	}
}

func TestWrite_CustomPrefix(t *testing.T) {
	var buf bytes.Buffer
	r := envwhistle.NewReporterWithPrefix(&buf, "[WARN]")
	findings := []envwhistle.Finding{
		{Key: "API_TOKEN", Severity: envwhistle.SeverityHigh, Reason: "token"},
	}
	_, _ = r.Write(findings)
	if !strings.HasPrefix(buf.String(), "[WARN]") {
		t.Errorf("expected custom prefix, got %q", buf.String())
	}
}

func TestSummary_Empty(t *testing.T) {
	s := envwhistle.Summary(nil)
	if s != "no findings" {
		t.Errorf("expected 'no findings', got %q", s)
	}
}

func TestSummary_MixedSeverities(t *testing.T) {
	findings := []envwhistle.Finding{
		{Severity: envwhistle.SeverityHigh},
		{Severity: envwhistle.SeverityHigh},
		{Severity: envwhistle.SeverityMedium},
		{Severity: envwhistle.SeverityLow},
	}
	s := envwhistle.Summary(findings)
	if !strings.Contains(s, "2 high") {
		t.Errorf("expected '2 high' in summary, got %q", s)
	}
	if !strings.Contains(s, "1 medium") {
		t.Errorf("expected '1 medium' in summary, got %q", s)
	}
	if !strings.Contains(s, "1 low") {
		t.Errorf("expected '1 low' in summary, got %q", s)
	}
}
