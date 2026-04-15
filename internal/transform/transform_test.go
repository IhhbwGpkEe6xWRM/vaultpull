package transform_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/transform"
)

func TestApply_EmptyPipeline(t *testing.T) {
	p := transform.New()
	src := map[string]string{"key": "value"}
	out := p.Apply(src)
	if out["key"] != "value" {
		t.Fatalf("expected value to pass through unchanged, got %q", out["key"])
	}
}

func TestUppercaseKeys(t *testing.T) {
	p := transform.New(transform.UppercaseKeys())
	out := p.Apply(map[string]string{"db_host": "localhost"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected DB_HOST key after uppercase transform")
	}
}

func TestPrefixKeys(t *testing.T) {
	p := transform.New(transform.PrefixKeys("APP_"))
	out := p.Apply(map[string]string{"SECRET": "abc"})
	if _, ok := out["APP_SECRET"]; !ok {
		t.Fatal("expected APP_SECRET after prefix transform")
	}
}

func TestPrefixKeys_EmptyPrefix(t *testing.T) {
	p := transform.New(transform.PrefixKeys(""))
	out := p.Apply(map[string]string{"KEY": "val"})
	if _, ok := out["KEY"]; !ok {
		t.Fatal("key should be unchanged with empty prefix")
	}
}

func TestTrimSpace(t *testing.T) {
	p := transform.New(transform.TrimSpace())
	out := p.Apply(map[string]string{"  KEY  ": "  value  "})
	if v, ok := out["KEY"]; !ok || v != "value" {
		t.Fatalf("expected trimmed key/value, got key present=%v value=%q", ok, v)
	}
}

func TestReplaceHyphens(t *testing.T) {
	p := transform.New(transform.ReplaceHyphens())
	out := p.Apply(map[string]string{"my-secret-key": "val"})
	if _, ok := out["my_secret_key"]; !ok {
		t.Fatal("expected hyphens replaced with underscores")
	}
}

func TestDropNonPrintable(t *testing.T) {
	p := transform.New(transform.DropNonPrintable())
	out := p.Apply(map[string]string{"KEY": "val\x00ue"})
	if out["KEY"] != "value" {
		t.Fatalf("expected non-printable chars removed, got %q", out["KEY"])
	}
}

func TestPipeline_ChainedTransformers(t *testing.T) {
	p := transform.New(
		transform.TrimSpace(),
		transform.ReplaceHyphens(),
		transform.UppercaseKeys(),
		transform.PrefixKeys("VAULT_"),
	)
	out := p.Apply(map[string]string{" db-password ": " s3cr3t "})
	v, ok := out["VAULT_DB_PASSWORD"]
	if !ok {
		t.Fatal("expected VAULT_DB_PASSWORD key")
	}
	if v != "s3cr3t" {
		t.Fatalf("expected trimmed value, got %q", v)
	}
}

func TestPipeline_DropsEmptyKey(t *testing.T) {
	// A transformer that blanks the key for a specific entry.
	dropFoo := func(k, v string) (string, string) {
		if k == "FOO" {
			return "", v
		}
		return k, v
	}
	p := transform.New(transform.Transformer(dropFoo))
	out := p.Apply(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if _, ok := out["FOO"]; ok {
		t.Fatal("expected FOO to be dropped")
	}
	if _, ok := out["BAZ"]; !ok {
		t.Fatal("expected BAZ to survive")
	}
}
