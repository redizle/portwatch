package watchlist_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func writeJSON(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "watchlist.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeJSON(t, `[
		{"port": 22,   "label": "ssh",  "priority": 3},
		{"port": 8080, "label": "dev",  "priority": 1}
	]`)

	w := watchlist.New()
	if err := w.LoadFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Contains(22) || !w.Contains(8080) {
		t.Error("expected ports 22 and 8080 to be loaded")
	}
	e, _ := w.Get(22)
	if e.Label != "ssh" {
		t.Errorf("expected label 'ssh', got %q", e.Label)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	w := watchlist.New()
	if err := w.LoadFile("/nonexistent/path.json"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	path := writeJSON(t, `not json at all`)
	w := watchlist.New()
	if err := w.LoadFile(path); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadFile_OutOfRangePriority_DefaultsToNormal(t *testing.T) {
	path := writeJSON(t, `[{"port": 3000, "label": "app", "priority": 99}]`)
	w := watchlist.New()
	if err := w.LoadFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := w.Get(3000)
	if !ok {
		t.Fatal("expected port 3000")
	}
	if e.Priority != watchlist.PriorityNormal {
		t.Errorf("expected PriorityNormal for out-of-range value, got %d", e.Priority)
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	path := writeJSON(t, `[{"port": 0, "label": "bad", "priority": 1}]`)
	w := watchlist.New()
	if err := w.LoadFile(path); err == nil {
		t.Error("expected error for invalid port")
	}
}
