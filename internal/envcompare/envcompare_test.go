package envcompare_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envcompare"
)

func TestCompare_AllMatch(t *testing.T) {
	c := envcompare.New(false)
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	r := c.Compare(left, right)
	if !r.Matches() {
		t.Fatal("expected all keys to match")
	}
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestCompare_Mismatch(t *testing.T) {
	c := envcompare.New(false)
	left := map[string]string{"KEY": "old"}
	right := map[string]string{"KEY": "new"}
	r := c.Compare(left, right)
	if r.Matches() {
		t.Fatal("expected mismatch")
	}
	if r.Entries[0].Status != envcompare.StatusMismatch {
		t.Fatalf("expected StatusMismatch, got %s", r.Entries[0].Status)
	}
}

func TestCompare_LeftOnly(t *testing.T) {
	c := envcompare.New(false)
	left := map[string]string{"ONLY_LEFT": "v"}
	right := map[string]string{}
	r := c.Compare(left, right)
	if r.Entries[0].Status != envcompare.StatusLeftOnly {
		t.Fatalf("expected StatusLeftOnly, got %s", r.Entries[0].Status)
	}
}

func TestCompare_RightOnly(t *testing.T) {
	c := envcompare.New(false)
	left := map[string]string{}
	right := map[string]string{"ONLY_RIGHT": "v"}
	r := c.Compare(left, right)
	if r.Entries[0].Status != envcompare.StatusRightOnly {
		t.Fatalf("expected StatusRightOnly, got %s", r.Entries[0].Status)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	c := envcompare.New(false)
	left := map[string]string{"Z": "1", "A": "2", "M": "3"}
	right := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := c.Compare(left, right)
	keys := []string{r.Entries[0].Key, r.Entries[1].Key, r.Entries[2].Key}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("expected sorted keys, got %v", keys)
	}
}

func TestCompare_MaskValues(t *testing.T) {
	c := envcompare.New(true)
	left := map[string]string{"SECRET": "plaintext"}
	right := map[string]string{"SECRET": "other"}
	r := c.Compare(left, right)
	if r.Entries[0].Left == "plaintext" {
		t.Fatal("expected value to be masked")
	}
	if r.Entries[0].Left != "<redacted>" {
		t.Fatalf("unexpected masked value: %s", r.Entries[0].Left)
	}
}

func TestCompare_MaskEmptyValue(t *testing.T) {
	c := envcompare.New(true)
	left := map[string]string{"K": ""}
	right := map[string]string{"K": ""}
	r := c.Compare(left, right)
	if r.Entries[0].Left != "<empty>" {
		t.Fatalf("expected <empty>, got %q", r.Entries[0].Left)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	c := envcompare.New(false)
	r := c.Compare(map[string]string{}, map[string]string{})
	if len(r.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(r.Entries))
	}
	if !r.Matches() {
		t.Fatal("empty maps should match")
	}
}
