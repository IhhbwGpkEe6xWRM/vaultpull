package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/"+path {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestResolvePath_NoNamespace(t *testing.T) {
	c := &Client{namespace: ""}
	got := c.resolvePath("secret/myapp")
	if got != "secret/myapp" {
		t.Errorf("expected %q, got %q", "secret/myapp", got)
	}
}

func TestResolvePath_WithNamespace(t *testing.T) {
	c := &Client{namespace: "team-a"}
	got := c.resolvePath("secret/myapp")
	want := "team-a/secret/myapp"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestResolvePath_StripsLeadingSlash(t *testing.T) {
	c := &Client{namespace: "team-a"}
	got := c.resolvePath("/secret/myapp")
	want := "team-a/secret/myapp"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestReadSecrets_KVv1(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{"DB_PASS": "secret123", "API_KEY": "abc"},
	}
	server := newTestServer(t, "secret/myapp", payload)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}

	if secrets["DB_PASS"] != "secret123" {
		t.Errorf("expected DB_PASS=secret123, got %q", secrets["DB_PASS"])
	}
	if secrets["API_KEY"] != "abc" {
		t.Errorf("expected API_KEY=abc, got %q", secrets["API_KEY"])
	}
}

func TestReadSecrets_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestToStringMap_SkipsNonString(t *testing.T) {
	in := map[string]interface{}{
		"KEY": "value",
		"NUM": 42,
		"BOOL": true,
	}
	out := toStringMap(in)
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", out["KEY"])
	}
}
