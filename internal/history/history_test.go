package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRecord_AddsEntry(t *testing.T) {
	h := New("", 100)
	h.Record(8080, "open")
	entries := h.GetAll()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Status != "open" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestRecord_MaxEntries(t *testing.T) {
	h := New("", 5)
	for i := 0; i < 10; i++ {
		h.Record(i, "open")
	}
	entries := h.GetAll()
	if len(entries) != 5 {
		t.Fatalf("expected 5 entries after cap, got %d", len(entries))
	}
	// should have kept the last 5
	if entries[0].Port != 5 {
		t.Errorf("expected oldest kept port to be 5, got %d", entries[0].Port)
	}
}

func TestGetByPort_Filters(t *testing.T) {
	h := New("", 100)
	h.Record(8080, "open")
	h.Record(9090, "open")
	h.Record(8080, "closed")

	results := h.GetByPort(8080)
	if len(results) != 2 {
		t.Fatalf("expected 2 entries for port 8080, got %d", len(results))
	}
}

func TestGetByPort_Missing(t *testing.T) {
	h := New("", 100)
	results := h.GetByPort(1234)
	if results != nil {
		t.Errorf("expected nil for missing port, got %v", results)
	}
}

func TestPersist_AndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h1 := New(path, 100)
	h1.Record(3000, "open")
	h1.Record(3001, "closed")

	// load into a second instance
	h2 := New(path, 100)
	entries := h2.GetAll()
	if len(entries) != 2 {
		t.Fatalf("expected 2 loaded entries, got %d", len(entries))
	}
}

func TestPersist_FileContents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := New(path, 100)
	h.Record(443, "open")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read persisted file: %v", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON in persisted file: %v", err)
	}
	if len(entries) != 1 || entries[0].Port != 443 {
		t.Errorf("unexpected persisted content: %+v", entries)
	}
}

func TestNew_DefaultMaxEntries(t *testing.T) {
	h := New("", 0)
	if h.maxEntries != 1000 {
		t.Errorf("expected default maxEntries 1000, got %d", h.maxEntries)
	}
}
