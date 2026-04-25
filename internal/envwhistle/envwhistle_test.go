package envwhistle_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envwhistle"
)

func TestScan_EmptyMap_NoFindings(t *testing.T) {
	d := envwhistle.New()
	findings := d.Scan(map[string]string{})
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestScan_HighSeverity_PasswordKey(t *testing.T) {
	d := envwhistle.New()
	findings := d.Scan(map[string]string{"DB_PASSWORD": "secret123"})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != envwhistle.SeverityHigh {
		t.Errorf("expected high severity, got %s", findings[0].Severity)
	}
}

func TestScan_MediumSeverity_AccessKey(t *testing.T) {
	d := envwhistle.New()
	findings := d.Scan(map[string]string{"AWS_ACCESS_KEY": "AKIAIOSFODNN7EXAMPLE"})
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	for _, f := range findings {
		if f.Key == "AWS_ACCESS_KEY" && f.Severity == envwhistle.SeverityMedium {
			return
		}
	}
	t.Error("expected medium severity finding for AWS_ACCESS_KEY")
}

func TestScan_NonSensitiveKey_NoFinding(t *testing.T) {
	d := envwhistle.New()
	findings := d.Scan(map[string]string{"APP_NAME": "vaultpull", "LOG_LEVEL": "info"})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-sensitive keys, got %d", len(findings))
	}
}

func TestScan_SortedByKey(t *testing.T) {
	d := envwhistle.New()
	findings := d.Scan(map[string]string{
		"Z_TOKEN": "t1",
		"A_SECRET": "s1",
		"M_PASSWORD": "p1",
	})
	if len(findings) < 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
	for i := 1; i < len(findings); i++ {
		if findings[i].Key < findings[i-1].Key {
			t.Errorf("findings not sorted: %s before %s", findings[i-1].Key, findings[i].Key)
		}
	}
}

func TestHasHigh_True(t *testing.T) {
	findings := []envwhistle.Finding{
		{Key: "X", Severity: envwhistle.SeverityLow},
		{Key: "Y", Severity: envwhistle.SeverityHigh},
	}
	if !envwhistle.HasHigh(findings) {
		t.Error("expected HasHigh to return true")
	}
}

func TestHasHigh_False(t *testing.T) {
	findings := []envwhistle.Finding{
		{Key: "X", Severity: envwhistle.SeverityLow},
		{Key: "Y", Severity: envwhistle.SeverityMedium},
	}
	if envwhistle.HasHigh(findings) {
		t.Error("expected HasHigh to return false")
	}
}

func TestNewWithRules_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envwhistle.NewWithRules([]string{"[invalid"}, envwhistle.SeverityHigh, "bad")
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestNewWithRules_CustomPattern_Matches(t *testing.T) {
	d, err := envwhistle.NewWithRules([]string{`(?i)my_custom`}, envwhistle.SeverityMedium, "custom rule")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	findings := d.Scan(map[string]string{"MY_CUSTOM_KEY": "value"})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Reason != "custom rule" {
		t.Errorf("unexpected reason: %s", findings[0].Reason)
	}
}
