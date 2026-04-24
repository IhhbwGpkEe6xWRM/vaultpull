package envsplit_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envsplit"
)

func TestNew_ValidRules(t *testing.T) {
	_, err := envsplit.New([]envsplit.Rule{
		{Name: "db", Prefix: "DB_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyName_ReturnsError(t *testing.T) {
	_, err := envsplit.New([]envsplit.Rule{
		{Name: "", Prefix: "DB_"},
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNew_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := envsplit.New([]envsplit.Rule{
		{Name: "db", Prefix: ""},
	})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestSplit_MatchesByPrefix(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Name: "db", Prefix: "DB_"},
		{Name: "api", Prefix: "API_"},
	})
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"API_KEY":     "secret",
		"OTHER_VALUE": "misc",
	}
	res := s.Split(secrets)

	if res.Groups["db"]["HOST"] != "localhost" {
		t.Errorf("expected db.HOST=localhost, got %q", res.Groups["db"]["HOST"])
	}
	if res.Groups["db"]["PORT"] != "5432" {
		t.Errorf("expected db.PORT=5432, got %q", res.Groups["db"]["PORT"])
	}
	if res.Groups["api"]["KEY"] != "secret" {
		t.Errorf("expected api.KEY=secret, got %q", res.Groups["api"]["KEY"])
	}
	if _, ok := res.Remainder["OTHER_VALUE"]; !ok {
		t.Error("expected OTHER_VALUE in remainder")
	}
}

func TestSplit_FirstMatchWins(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Name: "first", Prefix: "FOO_"},
		{Name: "second", Prefix: "FOO_BAR_"},
	})
	res := s.Split(map[string]string{"FOO_BAR_KEY": "val"})

	if _, ok := res.Groups["first"]["BAR_KEY"]; !ok {
		t.Error("expected first rule to win")
	}
	if len(res.Groups["second"]) != 0 {
		t.Error("second rule should not have matched")
	}
}

func TestSplit_CaseInsensitivePrefix(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Name: "svc", Prefix: "svc_"},
	})
	res := s.Split(map[string]string{"SVC_URL": "http://example.com"})

	if res.Groups["svc"]["URL"] != "http://example.com" {
		t.Errorf("expected case-insensitive match, got groups: %v", res.Groups)
	}
}

func TestSplit_EmptySecrets(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Name: "db", Prefix: "DB_"},
	})
	res := s.Split(map[string]string{})

	if len(res.Groups["db"]) != 0 {
		t.Error("expected empty db group")
	}
	if len(res.Remainder) != 0 {
		t.Error("expected empty remainder")
	}
}

func TestSplit_NoRules_AllRemainder(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{})
	res := s.Split(map[string]string{"FOO": "bar", "BAZ": "qux"})

	if len(res.Remainder) != 2 {
		t.Errorf("expected 2 remainder keys, got %d", len(res.Remainder))
	}
}
