package pagesize_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/pagesize"
)

func TestNew_ValidConfig(t *testing.T) {
	_, err := pagesize.New(pagesize.Config{PageSize: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_BelowMin(t *testing.T) {
	_, err := pagesize.New(pagesize.Config{PageSize: 0})
	if err == nil {
		t.Fatal("expected error for page size 0")
	}
}

func TestNew_AboveMax(t *testing.T) {
	_, err := pagesize.New(pagesize.Config{PageSize: 1001})
	if err == nil {
		t.Fatal("expected error for page size 1001")
	}
}

func TestPages_Empty(t *testing.T) {
	p, _ := pagesize.New(pagesize.DefaultConfig())
	if got := p.Pages(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestPages_SinglePage(t *testing.T) {
	p, _ := pagesize.New(pagesize.Config{PageSize: 10})
	keys := []string{"a", "b", "c"}
	pages := p.Pages(keys)
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if len(pages[0]) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(pages[0]))
	}
}

func TestPages_MultiplePages(t *testing.T) {
	p, _ := pagesize.New(pagesize.Config{PageSize: 3})
	keys := []string{"a", "b", "c", "d", "e"}
	pages := p.Pages(keys)
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
	if len(pages[0]) != 3 {
		t.Errorf("page 0: expected 3, got %d", len(pages[0]))
	}
	if len(pages[1]) != 2 {
		t.Errorf("page 1: expected 2, got %d", len(pages[1]))
	}
}

func TestTotalPages(t *testing.T) {
	p, _ := pagesize.New(pagesize.Config{PageSize: 10})
	cases := []struct {
		total, want int
	}{
		{0, 0},
		{1, 1},
		{10, 1},
		{11, 2},
		{25, 3},
	}
	for _, c := range cases {
		if got := p.TotalPages(c.total); got != c.want {
			t.Errorf("TotalPages(%d) = %d, want %d", c.total, got, c.want)
		}
	}
}
