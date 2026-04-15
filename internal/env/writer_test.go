package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_Write_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	w := NewWriter(p)
	err := w.Write(map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}
}

func TestWriter_Write_ContentsCorrect(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	w := NewWriter(p)
	err := w.Write(map[string]string{
		"db-password": "secret",
		"API_KEY":     "abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "API_KEY=abc123") {
		t.Errorf("expected API_KEY=abc123 in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PASSWORD=secret") {
		t.Errorf("expected DB_PASSWORD=secret in output, got:\n%s", content)
	}
}

func TestWriter_Write_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	w := NewWriter(p)
	err := w.Write(map[string]string{"MSG": "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWriter_Write_SortedKeys(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	w := NewWriter(p)
	err := w.Write(map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line to be ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line to be ZEBRA, got %s", lines[2])
	}
}

func TestWriter_Write_InvalidPath(t *testing.T) {
	w := NewWriter("/nonexistent/dir/.env")
	err := w.Write(map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
