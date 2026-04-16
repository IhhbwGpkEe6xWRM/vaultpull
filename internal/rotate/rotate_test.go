package rotate

import (
	"testing"
)

func TestDetect_NoChanges(t *testing.T) {
	d := NewDetector("secret/app")
	prev := map[string]string{"KEY": "value"}
	curr := map[string]string{"KEY": "value"}
	if events := d.Detect(prev, curr); len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestDetect_RotatedKey(t *testing.T) {
	d := NewDetector("secret/app")
	prev := map[string]string{"DB_PASS": "old", "API_KEY": "same"}
	curr := map[string]string{"DB_PASS": "new", "API_KEY": "same"}
	events := d.Detect(prev, curr)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %s", events[0].Key)
	}
	if events[0].Path != "secret/app" {
		t.Errorf("unexpected path %s", events[0].Path)
	}
}

func TestDetect_NewKeyNotRotation(t *testing.T) {
	d := NewDetector("secret/app")
	prev := map[string]string{}
	curr := map[string]string{"NEW_KEY": "val"}
	if events := d.Detect(prev, curr); len(events) != 0 {
		t.Fatalf("new keys should not be treated as rotations, got %d events", len(events))
	}
}

func TestDetect_MultipleRotations(t *testing.T) {
	d := NewDetector("secret/svc")
	prev := map[string]string{"A": "1", "B": "2", "C": "3"}
	curr := map[string]string{"A": "x", "B": "2", "C": "y"}
	events := d.Detect(prev, curr)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestSummary_NoEvents(t *testing.T) {
	if s := Summary(nil); s != "no rotated secrets detected" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithEvents(t *testing.T) {
	events := []Event{{Key: "K", Path: "p"}}
	if s := Summary(events); s != "1 secret(s) rotated" {
		t.Errorf("unexpected summary: %s", s)
	}
}
