package envstamp_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envstamp"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

var epoch = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestApply_NoOptions_ReturnsOriginal(t *testing.T) {
	s := envstamp.New()
	in := map[string]string{"FOO": "bar"}
	out := s.Apply(in)
	if out["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %q", out["FOO"])
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s := envstamp.New(envstamp.WithVersion("1.2.3"))
	in := map[string]string{"A": "1"}
	_ = s.Apply(in)
	if _, ok := in["BUILD_VERSION"]; ok {
		t.Fatal("Apply mutated input map")
	}
}

func TestApply_Version(t *testing.T) {
	s := envstamp.New(envstamp.WithVersion("v0.9.1"))
	out := s.Apply(map[string]string{})
	if out["BUILD_VERSION"] != "v0.9.1" {
		t.Fatalf("unexpected BUILD_VERSION: %q", out["BUILD_VERSION"])
	}
}

func TestApply_Timestamp(t *testing.T) {
	s := envstamp.New(envstamp.WithTimestamp(fixedClock(epoch)))
	out := s.Apply(map[string]string{})
	if out["STAMP_TIMESTAMP"] != "2024-06-01T12:00:00Z" {
		t.Fatalf("unexpected STAMP_TIMESTAMP: %q", out["STAMP_TIMESTAMP"])
	}
}

func TestApply_Hostname(t *testing.T) {
	s := envstamp.New(envstamp.WithHostname())
	out := s.Apply(map[string]string{})
	if out["STAMP_HOSTNAME"] == "" {
		t.Fatal("expected STAMP_HOSTNAME to be set")
	}
}

func TestApply_WithPrefix(t *testing.T) {
	s := envstamp.New(
		envstamp.WithPrefix("APP"),
		envstamp.WithVersion("2.0.0"),
		envstamp.WithTimestamp(fixedClock(epoch)),
	)
	out := s.Apply(map[string]string{})
	if _, ok := out["APP_BUILD_VERSION"]; !ok {
		t.Fatal("expected APP_BUILD_VERSION key")
	}
	if _, ok := out["APP_STAMP_TIMESTAMP"]; !ok {
		t.Fatal("expected APP_STAMP_TIMESTAMP key")
	}
	if _, ok := out["BUILD_VERSION"]; ok {
		t.Fatal("unprefixed BUILD_VERSION should not exist")
	}
}

func TestKeys_ReturnsConfiguredNames(t *testing.T) {
	s := envstamp.New(
		envstamp.WithVersion("1.0"),
		envstamp.WithTimestamp(fixedClock(epoch)),
		envstamp.WithHostname(),
	)
	keys := s.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(keys), keys)
	}
	joined := strings.Join(keys, ",")
	for _, want := range []string{"BUILD_VERSION", "STAMP_TIMESTAMP", "STAMP_HOSTNAME"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected key %q in Keys()", want)
		}
	}
}

func TestKeys_EmptyWhenNoOptions(t *testing.T) {
	s := envstamp.New()
	if len(s.Keys()) != 0 {
		t.Fatal("expected no keys")
	}
}
