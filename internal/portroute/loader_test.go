package portroute_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/user/portwatch/internal/portroute"
)

func writeRouteJSON(t *testing.T, data any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portroute-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatalf("encode: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeRouteJSON(t, []map[string]any{
		{"port": 80, "route": "/web"},
		{"port": 443, "route": "/secure"},
	})
	r := portroute.New()
	if err := portroute.LoadFile(path, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 2 {
		t.Fatalf("expected 2 routes, got %d", r.Len())
	}
	got, ok := r.Get(80)
	if !ok || got != "/web" {
		t.Errorf("expected /web for port 80, got %q", got)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	r := portroute.New()
	if err := portroute.LoadFile("/nonexistent/file.json", r); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "portroute-bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	r := portroute.New()
	if err := portroute.LoadFile(f.Name(), r); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeRouteJSON(t, []map[string]any{
		{"port": 0, "route": "/bad"},
	})
	r := portroute.New()
	if err := portroute.LoadFile(path, r); err == nil {
		t.Fatal("expected error for invalid port")
	}
}
