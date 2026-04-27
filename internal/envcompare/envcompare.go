// Package envcompare provides side-by-side comparison of two secret maps,
// producing structured results suitable for display or automated decision-making.
package envcompare

import "sort"

// Status indicates the comparison result for a single key.
type Status string

const (
	StatusMatch    Status = "match"
	StatusMismatch Status = "mismatch"
	StatusLeftOnly Status = "left_only"
	StatusRightOnly Status = "right_only"
)

// Entry holds the comparison result for one key.
type Entry struct {
	Key    string
	Status Status
	Left   string
	Right  string
}

// Result is the full comparison output.
type Result struct {
	Entries []Entry
}

// Matches returns true when every key matches across both maps.
func (r Result) Matches() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return false
		}
	}
	return true
}

// Comparer compares two string maps.
type Comparer struct {
	maskValues bool
}

// New returns a Comparer. When maskValues is true, values are replaced with
// redacted placeholders in the output so secrets are not exposed in logs.
func New(maskValues bool) *Comparer {
	return &Comparer{maskValues: maskValues}
}

// Compare performs the side-by-side comparison of left and right.
func (c *Comparer) Compare(left, right map[string]string) Result {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		lv, lok := left[k]
		rv, rok := right[k]

		var status Status
		switch {
		case lok && !rok:
			status = StatusLeftOnly
		case !lok && rok:
			status = StatusRightOnly
		case lv == rv:
			status = StatusMatch
		default:
			status = StatusMismatch
		}

		if c.maskValues {
			lv = mask(lv, lok)
			rv = mask(rv, rok)
		}

		entries = append(entries, Entry{Key: k, Status: status, Left: lv, Right: rv})
	}
	return Result{Entries: entries}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

func mask(v string, present bool) string {
	if !present {
		return ""
	}
	if len(v) == 0 {
		return "<empty>"
	}
	return "<redacted>"
}
