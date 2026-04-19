package tokenize

import (
	"testing"
)

func TestParse_SimplePath(t *testing.T) {
	tok := New().Parse("secret/app/db")
	if len(tok.Parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(tok.Parts))
	}
	if tok.Parts[0] != "secret" || tok.Parts[1] != "app" || tok.Parts[2] != "db" {
		t.Errorf("unexpected parts: %v", tok.Parts)
	}
}

func TestParse_TrimsLeadingTrailingSlashes(t *testing.T) {
	tok := New().Parse("/secret/app/")
	if len(tok.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(tok.Parts))
	}
}

func TestParse_EmptyPath(t *testing.T) {
	tok := New().Parse("")
	if len(tok.Parts) != 0 {
		t.Errorf("expected empty parts for empty path")
	}
}

func TestParse_SlashOnly(t *testing.T) {
	tok := New().Parse("/")
	if len(tok.Parts) != 0 {
		t.Errorf("expected empty parts for slash-only path")
	}
}

func TestJoin_ReconstructsPath(t *testing.T) {
	tz := New()
	tok := tz.Parse("secret/app/db")
	joined := tz.Join(tok)
	if joined != "secret/app/db" {
		t.Errorf("expected 'secret/app/db', got %q", joined)
	}
}

func TestDepth_CountsSegments(t *testing.T) {
	tz := New()
	tok := tz.Parse("a/b/c/d")
	if tz.Depth(tok) != 4 {
		t.Errorf("expected depth 4, got %d", tz.Depth(tok))
	}
}

func TestParent_ReturnsPart(t *testing.T) {
	tz := New()
	tok := tz.Parse("secret/app/db")
	parent, err := tz.Parent(tok)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz.Join(parent) != "secret/app" {
		t.Errorf("expected 'secret/app', got %q", tz.Join(parent))
	}
}

func TestParent_ErrorOnRoot(t *testing.T) {
	tz := New()
	tok := tz.Parse("")
	_, err := tz.Parent(tok)
	if err == nil {
		t.Error("expected error for root token, got nil")
	}
}

func TestCustomSeparator(t *testing.T) {
	tz := NewWithSeparator(".")
	tok := tz.Parse("secret.app.db")
	if len(tok.Parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(tok.Parts))
	}
	if tz.Join(tok) != "secret.app.db" {
		t.Errorf("unexpected join result: %q", tz.Join(tok))
	}
}
