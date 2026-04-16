package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/prompt"
)

func newTestConfirmer(input string) (*prompt.Confirmer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	c := prompt.NewWithReadWriter(strings.NewReader(input), out)
	return c, out
}

func TestConfirm_YesLower(t *testing.T) {
	c, _ := newTestConfirmer("y\n")
	ok, err := c.Confirm("overwrite?")
	if err != nil || !ok {
		t.Fatalf("expected true, got ok=%v err=%v", ok, err)
	}
}

func TestConfirm_YesFull(t *testing.T) {
	c, _ := newTestConfirmer("yes\n")
	ok, err := c.Confirm("overwrite?")
	if err != nil || !ok {
		t.Fatalf("expected true, got ok=%v err=%v", ok, err)
	}
}

func TestConfirm_YesUppercase(t *testing.T) {
	c, _ := newTestConfirmer("YES\n")
	ok, err := c.Confirm("overwrite?")
	if err != nil || !ok {
		t.Fatalf("expected true for uppercase YES")
	}
}

func TestConfirm_No(t *testing.T) {
	c, _ := newTestConfirmer("n\n")
	ok, err := c.Confirm("overwrite?")
	if err != nil || ok {
		t.Fatalf("expected false for 'n'")
	}
}

func TestConfirm_EmptyInput(t *testing.T) {
	c, _ := newTestConfirmer("\n")
	ok, err := c.Confirm("overwrite?")
	if err != nil || ok {
		t.Fatalf("expected false for empty input")
	}
}

func TestConfirm_EOF(t *testing.T) {
	c, _ := newTestConfirmer("")
	ok, err := c.Confirm("overwrite?")
	if err != nil || ok {
		t.Fatalf("expected false on EOF, got ok=%v err=%v", ok, err)
	}
}

func TestConfirm_PrintsPrompt(t *testing.T) {
	c, out := newTestConfirmer("y\n")
	_, _ = c.Confirm("delete secrets?")
	if !strings.Contains(out.String(), "delete secrets?") {
		t.Fatalf("prompt text not written to output: %q", out.String())
	}
}

func TestConfirm_PrintsOptions(t *testing.T) {
	c, out := newTestConfirmer("n\n")
	_, _ = c.Confirm("continue?")
	if !strings.Contains(out.String(), "[y/N]") {
		t.Fatalf("expected [y/N] hint in output: %q", out.String())
	}
}
