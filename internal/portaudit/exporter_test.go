package portaudit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSON_ValidOutput(t *testing.T) {
	l := New(0)
	l.Record(80, ActionOpened, "http")
	ex := NewExporter(l)
	var buf bytes.Buffer
	if err := ex.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 80 {
		t.Errorf("expected port 80, got %d", entries[0].Port)
	}
}

func TestWriteText_ContainsAction(t *testing.T) {
	l := New(0)
	l.Record(443, ActionAlerted, "spike")
	ex := NewExporter(l)
	var buf bytes.Buffer
	if err := ex.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alerted") {
		t.Errorf("expected 'alerted' in output, got: %s", buf.String())
	}
}

func TestWriteJSON_EmptyLog(t *testing.T) {
	l := New(0)
	ex := NewExporter(l)
	var buf bytes.Buffer
	if err := ex.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Fatal("expected non-empty JSON output")
	}
}
