package envimport_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envimport"
)

func writeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestParse_SimpleKeyValue(t *testing.T) {
	r := strings.NewReader("FOO=bar\nBAZ=qux\n")
	m, err := envimport.Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParse_SkipsComments(t *testing.T) {
	r := strings.NewReader("# comment\nKEY=val\n")
	m, err := envimport.Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 1 || m["KEY"] != "val" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParse_UnquotesDoubleQuotes(t *testing.T) {
	r := strings.NewReader(`SECRET="hello world"` + "\n")
	m, err := envimport.Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["SECRET"] != "hello world" {
		t.Fatalf("got %q", m["SECRET"])
	}
}

func TestParse_UnquotesSingleQuotes(t *testing.T) {
	r := strings.NewReader("TOKEN='abc123'\n")
	m, err := envimport.Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["TOKEN"] != "abc123" {
		t.Fatalf("got %q", m["TOKEN"])
	}
}

func TestParse_MissingEquals_ReturnsError(t *testing.T) {
	r := strings.NewReader("BADLINE\n")
	_, err := envimport.Parse(r)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_EmptyKey_ReturnsError(t *testing.T) {
	r := strings.NewReader("=value\n")
	_, err := envimport.Parse(r)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	i := envimport.New("/nonexistent/.env")
	_, err := i.Load()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_ReadsFile(t *testing.T) {
	p := writeEnv(t, "DB_URL=postgres://localhost\nPORT=5432\n")
	i := envimport.New(p)
	m, err := i.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DB_URL"] != "postgres://localhost" {
		t.Fatalf("unexpected value: %v", m)
	}
	if m["PORT"] != "5432" {
		t.Fatalf("unexpected value: %v", m)
	}
}
