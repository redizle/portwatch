package portalias

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "dev-server"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(8080)
	if !ok || v != "dev-server" {
		t.Fatalf("expected dev-server, got %q ok=%v", v, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "bad"); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_EmptyAlias(t *testing.T) {
	s := New()
	if err := s.Set(80, ""); err != ErrEmptyAlias {
		t.Fatalf("expected ErrEmptyAlias, got %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(443, "https")
	s.Remove(443)
	_, ok := s.Get(443)
	if ok {
		t.Fatal("expected removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(22, "ssh")
	_ = s.Set(80, "http")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestLoadFile_Valid(t *testing.T) {
	entries := []map[string]interface{}{
		{"port": 3306, "label": "mysql"},
		{"port": 5432, "label": "postgres"},
	}
	f, _ := os.CreateTemp(t.TempDir(), "aliases*.json")
	_ = json.NewEncoder(f).Encode(entries)
	f.Close()

	s := New()
	if err := s.LoadFile(f.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(3306)
	if !ok || v != "mysql" {
		t.Fatalf("expected mysql, got %q", v)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	s := New()
	if err := s.LoadFile("/no/such/file.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad*.json")
	_, _ = f.WriteString("not json")
	f.Close()
	s := New()
	if err := s.LoadFile(f.Name()); err == nil {
		t.Fatal("expected JSON decode error")
	}
}
