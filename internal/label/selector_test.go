package label_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/label"
)

func TestFilter_MatchingMeta(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "secret1",
		"API_KEY": "secret2",
	}
	meta := map[string]label.Set{
		"DB_PASS": {"tier": "db"},
		"API_KEY": {"tier": "api"},
	}
	sel := label.NewSelector(label.Set{"tier": "db"})
	out := sel.Filter(secrets, func(k string) label.Set { return meta[k] })
	if len(out) != 1 || out["DB_PASS"] != "secret1" {
		t.Fatalf("unexpected filter result: %v", out)
	}
}

func TestFilter_EmptyFilter_ReturnsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	sel := label.NewSelector(label.Set{})
	out := sel.Filter(secrets, func(k string) label.Set { return label.Set{} })
	if len(out) != 2 {
		t.Fatalf("expected all secrets, got %d", len(out))
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	sel := label.NewSelector(label.Set{"env": "prod"})
	out := sel.Filter(secrets, func(k string) label.Set { return label.Set{"env": "dev"} })
	if len(out) != 0 {
		t.Fatal("expected empty result")
	}
}

func TestParseFilter_ReturnsParsedSet(t *testing.T) {
	f := label.ParseFilter([]string{"env=prod", "region=eu"})
	if f["env"] != "prod" || f["region"] != "eu" {
		t.Fatalf("unexpected filter: %v", f)
	}
}
