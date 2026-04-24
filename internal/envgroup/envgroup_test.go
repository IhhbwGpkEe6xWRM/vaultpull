package envgroup_test

import (
	"sort"
	"testing"

	"github.com/your-org/vaultpull/internal/envgroup"
)

func groupByName(groups []envgroup.Group, name string) map[string]string {
	for _, g := range groups {
		if g.Name == name {
			return g.Values
		}
	}
	return nil
}

func groupNames(groups []envgroup.Group) []string {
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	sort.Strings(names)
	return names
}

func TestNew_ValidPairs(t *testing.T) {
	g, err := envgroup.New([]string{"db=DB", "cache=CACHE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil Grouper")
	}
}

func TestNew_MalformedPair(t *testing.T) {
	_, err := envgroup.New([]string{"NOEQUALSSIGN"})
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestNew_EmptyPrefix(t *testing.T) {
	_, err := envgroup.New([]string{"db="})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestSplit_MatchesByPrefix(t *testing.T) {
	g, _ := envgroup.New([]string{"db=DB"})
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "prod",
	}
	groups := g.Split(secrets)

	db := groupByName(groups, "db")
	if db["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", db["HOST"])
	}
	if db["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", db["PORT"])
	}
}

func TestSplit_UnmatchedGoesToDefault(t *testing.T) {
	g, _ := envgroup.New([]string{"db=DB"})
	secrets := map[string]string{
		"APP_ENV": "staging",
	}
	groups := g.Split(secrets)

	def := groupByName(groups, "")
	if def["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV in default group")
	}
}

func TestSplit_FirstMatchWins(t *testing.T) {
	g, _ := envgroup.New([]string{"a=DB", "b=DB_REPLICA"})
	secrets := map[string]string{
		"DB_REPLICA_HOST": "replica",
	}
	groups := g.Split(secrets)

	a := groupByName(groups, "a")
	if a == nil || a["REPLICA_HOST"] != "replica" {
		t.Errorf("expected first rule to win; got groups: %v", groupNames(groups))
	}
}

func TestSplit_EmptySecrets(t *testing.T) {
	g, _ := envgroup.New([]string{"db=DB"})
	groups := g.Split(map[string]string{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestSplit_CaseInsensitiveMatch(t *testing.T) {
	g, _ := envgroup.New([]string{"db=db"})
	secrets := map[string]string{"DB_NAME": "mydb"}
	groups := g.Split(secrets)
	db := groupByName(groups, "db")
	if db["NAME"] != "mydb" {
		t.Errorf("expected case-insensitive match; got %v", db)
	}
}
