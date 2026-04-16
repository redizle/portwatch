package digest

import (
	"strings"
	"testing"
	"time"
)

func namer(p int) string {
	names := map[int]string{80: "http", 443: "https", 22: "ssh"}
	if n, ok := names[p]; ok {
		return n
	}
	return "unknown"
}

func TestBuilder_Record_IncreasesLen(t *testing.T) {
	b := NewBuilder(time.Minute, namer)
	b.Record(80, "open")
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
}

func TestBuilder_Build_ReturnsEntries(t *testing.T) {
	b := NewBuilder(time.Minute, namer)
	b.Record(443, "open")
	b.Record(22, "closed")
	entries := b.Build()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestBuilder_Build_ClearsBuffer(t *testing.T) {
	b := NewBuilder(time.Minute, namer)
	b.Record(80, "open")
	b.Build()
	if b.Len() != 0 {
		t.Fatal("expected empty buffer after Build")
	}
}

func TestBuilder_NilNamer_UsesUnknown(t *testing.T) {
	b := NewBuilder(time.Minute, nil)
	b.Record(9999, "open")
	s := b.Summary()
	if !strings.Contains(s, "unknown") {
		t.Errorf("expected 'unknown' name, got: %s", s)
	}
}

func TestBuilder_Summary_ContainsPort(t *testing.T) {
	b := NewBuilder(time.Minute, namer)
	b.Record(443, "open")
	s := b.Summary()
	if !strings.Contains(s, "443") {
		t.Errorf("expected port 443 in summary: %s", s)
	}
}
