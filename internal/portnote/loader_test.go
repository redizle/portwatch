package portnote

import (
	"encoding/json"
	"os"
	"testing"
)

func writeNoteJSON(t *testing.T, entries []fileEntry) string {
	t.Helper()
	f, err := os.CreateTemp("", "portnote*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeNoteJSON(t, []fileEntry{
		{Port: 80, Text: "http traffic"},
		{Port: 443, Text: "tls"},
	})
	s := New()
	if err := LoadFile(s, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Errorf("expected 2 notes, got %d", s.Len())
	}
	n, ok := s.Get(80)
	if !ok || n.Text != "http traffic" {
		t.Errorf("unexpected note: %+v", n)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	s := New()
	if err := LoadFile(s, "/no/such/file.json"); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "portnote*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	s := New()
	if err := LoadFile(s, f.Name()); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeNoteJSON(t, []fileEntry{{Port: 0, Text: "bad"}})
	s := New()
	if err := LoadFile(s, path); err == nil {
		t.Fatal("expected error for invalid port")
	}
}
