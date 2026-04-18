package resolve_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/resolve"
)

func TestNew_InvalidMount(t *testing.T) {
	_, err := resolve.New("", resolve.KVv1)
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestNew_InvalidVersion(t *testing.T) {
	_, err := resolve.New("secret", resolve.MountVersion(99))
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestNew_TrimsMountSlashes(t *testing.T) {
	r, err := resolve.New("/secret/", resolve.KVv1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Mount() != "secret" {
		t.Errorf("expected 'secret', got %q", r.Mount())
	}
}

func TestDataPath_KVv1(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv1)
	got := r.DataPath("myapp/config")
	want := "secret/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDataPath_KVv2(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv2)
	got := r.DataPath("myapp/config")
	want := "secret/data/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDataPath_StripsLeadingSlash(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv2)
	got := r.DataPath("/myapp/config")
	want := "secret/data/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMetadataPath_KVv2(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv2)
	got := r.MetadataPath("myapp/config")
	want := "secret/metadata/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMetadataPath_KVv1_FallsBackToDataPath(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv1)
	got := r.MetadataPath("myapp/config")
	want := "secret/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestVersion_ReturnsConfigured(t *testing.T) {
	r, _ := resolve.New("secret", resolve.KVv2)
	if r.Version() != resolve.KVv2 {
		t.Errorf("expected KVv2, got %v", r.Version())
	}
}
