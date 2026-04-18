package environ_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/environ"
)

func TestLoad_NoPrefix_ReturnsAll(t *testing.T) {
	t.Setenv("SOME_VAR", "hello")
	l := environ.New("", false)
	got := l.Load()
	if got["SOME_VAR"] != "hello" {
		t.Errorf("expected SOME_VAR=hello, got %q", got["SOME_VAR"])
	}
}

func TestLoad_WithPrefix_FiltersKeys(t *testing.T) {
	t.Setenv("APP_DB_HOST", "localhost")
	t.Setenv("OTHER_KEY", "ignored")
	l := environ.New("APP", false)
	got := l.Load()
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if _, ok := got["OTHER_KEY"]; ok {
		t.Error("OTHER_KEY should not be present")
	}
}

func TestLoad_PrefixCaseInsensitive(t *testing.T) {
	t.Setenv("APP_SECRET", "s3cr3t")
	l := environ.New("app", false)
	got := l.Load()
	if got["SECRET"] != "s3cr3t" {
		t.Errorf("expected SECRET=s3cr3t, got %q", got["SECRET"])
	}
}

func TestMerge_SecretsWinByDefault(t *testing.T) {
	t.Setenv("APP_KEY", "from-env")
	l := environ.New("APP", false)
	secrets := map[string]string{"KEY": "from-vault"}
	out := l.Merge(secrets)
	if out["KEY"] != "from-vault" {
		t.Errorf("expected from-vault, got %q", out["KEY"])
	}
}

func TestMerge_OverrideTrue_EnvWins(t *testing.T) {
	t.Setenv("APP_KEY", "from-env")
	l := environ.New("APP", true)
	secrets := map[string]string{"KEY": "from-vault"}
	out := l.Merge(secrets)
	if out["KEY"] != "from-env" {
		t.Errorf("expected from-env, got %q", out["KEY"])
	}
}

func TestMerge_EnvOnlyKey_AlwaysIncluded(t *testing.T) {
	t.Setenv("APP_EXTRA", "bonus")
	l := environ.New("APP", false)
	out := l.Merge(map[string]string{})
	if out["EXTRA"] != "bonus" {
		t.Errorf("expected EXTRA=bonus, got %q", out["EXTRA"])
	}
}
