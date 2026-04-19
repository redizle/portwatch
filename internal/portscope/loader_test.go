package portscope

import (
	"encoding/json"
	"os"
	"testing"
)

func writeScopeJSON(t *testing.T, entries []fileEntry) string {
	t.Helper()
	f, err := os.CreateTemp("", "portscope-*.json")
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
	path := writeScopeJSON(t, []fileEntry{
		{Port: 80, Scope: "external"},
		{Port: 22, Scope: "internal"},
	})
	c := New(ScopeUnknown)
	if err := LoadFile(path, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Get(80); got != ScopeExternal {
		t.Errorf("expected external, got %s", got)
	}
	if got := c.Get(22); got != ScopeInternal {
		t.Errorf("expected internal, got %s", got)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	c := New(ScopeUnknown)
	if err := LoadFile("/nonexistent/file.json", c); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "portscope-bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	c := New(ScopeUnknown)
	if err := LoadFile(f.Name(), c); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeScopeJSON(t, []fileEntry{{Port: 0, Scope: "external"}})
	c := New(ScopeUnknown)
	if err := LoadFile(path, c); err == nil {
		t.Error("expected error for invalid port")
	}
}
