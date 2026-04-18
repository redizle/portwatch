package portcomment

import (
	"os"
	"path/filepath"
	"testing"
)

func writeJSON(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "comments*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeJSON(t, `[{"port":80,"comment":"http"},{"port":443,"comment":"https"}]`)
	s := New()
	if err := s.LoadFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c, ok := s.Get(80); !ok || c != "http" {
		t.Fatalf("expected http, got %q", c)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	s := New()
	err := s.LoadFile(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	path := writeJSON(t, `not json`)
	s := New()
	if err := s.LoadFile(path); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeJSON(t, `[{"port":0,"comment":"bad"}]`)
	s := New()
	if err := s.LoadFile(path); err == nil {
		t.Fatal("expected error for invalid port")
	}
}
