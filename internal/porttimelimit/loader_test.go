package porttimelimit

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTimelimitJSON(t *testing.T, entries []jsonEntry) string {
	t.Helper()
	f, err := os.CreateTemp("", "timelimit-*.json")
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
	entries := []jsonEntry{
		{Port: 8080, Start: "09:00", End: "17:00"},
		{Port: 443, Start: "08:00", End: "20:00"},
	}
	path := writeTimelimitJSON(t, entries)
	l, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, ok := l.Get(8080)
	if !ok {
		t.Fatal("expected entry for 8080")
	}
	if w.Start != 9*time.Hour || w.End != 17*time.Hour {
		t.Errorf("unexpected window: %+v", w)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/nonexistent/timelimit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	_, err := LoadFile(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadFile_InvalidPort(t *testing.T) {
	entries := []jsonEntry{{Port: 0, Start: "09:00", End: "17:00"}}
	path := writeTimelimitJSON(t, entries)
	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestLoadFile_BadClockFormat(t *testing.T) {
	entries := []jsonEntry{{Port: 80, Start: "9am", End: "17:00"}}
	path := writeTimelimitJSON(t, entries)
	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for bad clock format")
	}
}
