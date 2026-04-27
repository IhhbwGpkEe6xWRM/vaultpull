package envcompare_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envcompare"
)

func buildResult(left, right map[string]string, mask bool) envcompare.Result {
	c := envcompare.New(mask)
	return c.Compare(left, right)
}

func TestFormat_NoEntries(t *testing.T) {
	var sb strings.Builder
	r := buildResult(nil, nil, false)
	envcompare.Format(&sb, r, "vault", "local")
	if !strings.Contains(sb.String(), "no keys") {
		t.Fatalf("unexpected output: %s", sb.String())
	}
}

func TestFormat_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	r := buildResult(map[string]string{"X": "1"}, map[string]string{"X": "1"}, false)
	envcompare.Format(&sb, r, "vault", "local")
	out := sb.String()
	if !strings.Contains(out, "vault") {
		t.Fatal("expected left label in output")
	}
	if !strings.Contains(out, "local") {
		t.Fatal("expected right label in output")
	}
}

func TestFormat_ShowsMismatchSymbol(t *testing.T) {
	var sb strings.Builder
	r := buildResult(
		map[string]string{"DB_PASS": "old"},
		map[string]string{"DB_PASS": "new"},
		false,
	)
	envcompare.Format(&sb, r, "vault", "local")
	if !strings.Contains(sb.String(), "≠") {
		t.Fatal("expected mismatch symbol")
	}
}

func TestFormat_SummaryLine(t *testing.T) {
	var sb strings.Builder
	r := buildResult(
		map[string]string{"A": "1", "B": "x"},
		map[string]string{"A": "1", "B": "y"},
		false,
	)
	envcompare.Format(&sb, r, "left", "right")
	out := sb.String()
	if !strings.Contains(out, "1 match") {
		t.Fatalf("expected summary with 1 match, got:\n%s", out)
	}
	if !strings.Contains(out, "1 mismatch") {
		t.Fatalf("expected summary with 1 mismatch, got:\n%s", out)
	}
}

func TestFormat_TruncatesLongValues(t *testing.T) {
	var sb strings.Builder
	long := strings.Repeat("x", 50)
	r := buildResult(map[string]string{"K": long}, map[string]string{"K": long}, false)
	envcompare.Format(&sb, r, "l", "r")
	if strings.Contains(sb.String(), long) {
		t.Fatal("expected long value to be truncated")
	}
}
