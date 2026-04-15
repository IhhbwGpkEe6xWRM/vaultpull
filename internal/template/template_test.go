package template_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/template"
)

var sampleSecrets = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"API_KEY": "s3cr3t",
}

func TestNew_ValidTemplate(t *testing.T) {
	_, err := template.New(`{{ index . "DB_HOST" }}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := template.New(`{{ .Unclosed `)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRender_InterpolatesValues(t *testing.T) {
	r, err := template.New(`host={{ index . "DB_HOST" }} port={{ index . "DB_PORT" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.Render(sampleSecrets)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if !strings.Contains(out, "host=localhost") {
		t.Errorf("expected 'host=localhost' in output, got: %s", out)
	}
	if !strings.Contains(out, "port=5432") {
		t.Errorf("expected 'port=5432' in output, got: %s", out)
	}
}

func TestRender_UpperFuncMap(t *testing.T) {
	r, err := template.New(`{{ index . "DB_HOST" | upper }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.Render(sampleSecrets)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if out != "LOCALHOST" {
		t.Errorf("expected 'LOCALHOST', got: %s", out)
	}
}

func TestRenderToFile_WritesOutput(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "rendered.txt")

	r, err := template.New(`key={{ index . "API_KEY" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.RenderToFile(sampleSecrets, out); err != nil {
		t.Fatalf("RenderToFile error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if string(data) != "key=s3cr3t" {
		t.Errorf("unexpected file content: %s", string(data))
	}
}

func TestNewFromFile_ParsesTemplate(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "tpl.txt")
	if err := os.WriteFile(tmplPath, []byte(`db={{ index . "DB_HOST" }}`), 0644); err != nil {
		t.Fatalf("writing template file: %v", err)
	}

	r, err := template.NewFromFile(tmplPath)
	if err != nil {
		t.Fatalf("NewFromFile error: %v", err)
	}
	out, err := r.Render(sampleSecrets)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if out != "db=localhost" {
		t.Errorf("expected 'db=localhost', got: %s", out)
	}
}

func TestNewFromFile_MissingFile(t *testing.T) {
	_, err := template.NewFromFile("/nonexistent/path/tpl.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
