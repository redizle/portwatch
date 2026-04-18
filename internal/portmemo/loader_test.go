package portmemo_test

import (
	"encoding/json"
	"os"
	"testing"

	"portwatch/internal/portmemo"
)

func writeMemoJSON(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portmemo-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = json.NewEncoder(f).Encode(v)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeMemoJSON(t, []map[string]any{
		{"port": 80, "text": "http server", "ttl_seconds": 0},
		{"port": 443, "text": "https", "ttl_seconds": 0},
	})
	s := portmemo.New()
	if err := portmemo.LoadFile(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m, ok := s.Get(80); !ok || m.Text != "http server" {
		t.Error("expected memo for port 80")
	}
	if m, ok := s.Get(443); !ok || m.Text != "https" {
		t.Error("expected memo for port 443")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	s := portmemo.New()
	if err := portmemo.LoadFile("/no/such/file.json", s); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "portmemo-bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	s := portmemo.New()
	if err := portmemo.LoadFile(f.Name(), s); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeMemoJSON(t, []map[string]any{
		{"port": 99999, "text": "bad port", "ttl_seconds": 0},
	})
	s := portmemo.New()
	if err := portmemo.LoadFile(path, s); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}
