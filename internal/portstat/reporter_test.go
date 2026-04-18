package portstat

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoStats(t *testing.T) {
	tr := New()
	r := NewReporter(tr)
	var buf bytes.Buffer
	if err := r.Print(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no port statistics") {
		t.Errorf("expected no-stats message, got: %s", buf.String())
	}
}

func TestPrint_WithStats(t *testing.T) {
	tr := New()
	_ = tr.Record(80, "open")
	_ = tr.Record(80, "open")
	_ = tr.Record(443, "closed")
	r := NewReporter(tr)
	var buf bytes.Buffer
	if err := r.Print(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected header row")
	}
}

func TestPrint_SortedByPort(t *testing.T) {
	tr := New()
	_ = tr.Record(9000, "open")
	_ = tr.Record(22, "open")
	_ = tr.Record(80, "closed")
	r := NewReporter(tr)
	var buf bytes.Buffer
	_ = r.Print(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// lines[0] is header
	if !strings.HasPrefix(strings.TrimSpace(lines[1]), "22") {
		t.Errorf("expected first data row to be port 22, got: %s", lines[1])
	}
}
