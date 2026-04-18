package portannot

import (
	"encoding/json"
	"os"
	"testing"
)

func writeAnnotJSON(t *testing.T, data any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "annot*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadFile_Valid(t *testing.T) {
	payload := []map[string]any{
		{"port": 80, "annotations": map[string]string{"env": "prod", "team": "web"}},
	}
	path := writeAnnotJSON(t, payload)
	s := New()
	if err := LoadFile(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(80, "env")
	if !ok || v != "prod" {
		t.Fatalf("expected prod, got %q", v)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	s := New()
	if err := LoadFile("/no/such/file.json", s); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad*.json")
	f.WriteString("not json")
	f.Close()
	s := New()
	if err := LoadFile(f.Name(), s); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	payload := []map[string]any{
		{"port": 99999, "annotations": map[string]string{"k": "v"}},
	}
	path := writeAnnotJSON(t, payload)
	s := New()
	if err := LoadFile(path, s); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}
