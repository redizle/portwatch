package digest

import (
	"strings"
	"testing"
	"time"
)

func TestAdd_IncreasesLen(t *testing.T) {
	d := New(time.Minute)
	d.Add(Entry{Port: 80, Status: "open", Name: "http"})
	if d.Len() != 1 {
		t.Fatalf("expected 1, got %d", d.Len())
	}
}

func TestFlush_ReturnsAndClears(t *testing.T) {
	d := New(time.Minute)
	d.Add(Entry{Port: 443, Status: "open", Name: "https"})
	out := d.Flush()
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if d.Len() != 0 {
		t.Fatal("expected digest to be cleared after flush")
	}
}

func TestFlush_ExcludesExpiredEntries(t *testing.T) {
	d := New(time.Second)
	old := Entry{Port: 22, Status: "open", Name: "ssh", At: time.Now().Add(-2 * time.Second)}
	d.Add(old)
	out := d.Flush()
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestSummary_NoEntries(t *testing.T) {
	d := New(time.Minute)
	if d.Summary() != "no activity" {
		t.Fatal("expected 'no activity'")
	}
}

func TestSummary_WithEntries(t *testing.T) {
	d := New(time.Minute)
	d.Add(Entry{Port: 8080, Status: "open", Name: "http-alt"})
	s := d.Summary()
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", s)
	}
	if !strings.Contains(s, "http-alt") {
		t.Errorf("expected name in summary, got: %s", s)
	}
}

func TestAdd_SetsTimestamp(t *testing.T) {
	d := New(time.Minute)
	before := time.Now()
	d.Add(Entry{Port: 3306, Status: "open", Name: "mysql"})
	after := time.Now()
	e := d.entries[0]
	if e.At.Before(before) || e.At.After(after) {
		t.Error("timestamp not set correctly")
	}
}
