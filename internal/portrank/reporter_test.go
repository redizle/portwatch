package portrank

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoEntries(t *testing.T) {
	r := New()
	rp := NewReporter(r)
	var buf bytes.Buffer
	rp.Print(&buf)
	if !strings.Contains(buf.String(), "no entries") {
		t.Errorf("expected 'no entries' message, got: %s", buf.String())
	}
}

func TestPrint_WithEntries(t *testing.T) {
	r := New()
	_ = r.Add(80, 10)
	_ = r.Add(443, 3)
	rp := NewReporter(r)
	var buf bytes.Buffer
	rp.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
}

func TestPrint_SortedByScore(t *testing.T) {
	r := New()
	_ = r.Add(9000, 1)
	_ = r.Add(22, 50)
	_ = r.Add(80, 25)
	rp := NewReporter(r)
	var buf bytes.Buffer
	rp.Print(&buf)
	out := buf.String()
	idx22 := strings.Index(out, "22")
	idx80 := strings.Index(out, "80")
	idx9000 := strings.Index(out, "9000")
	if idx22 > idx80 || idx80 > idx9000 {
		t.Errorf("expected order 22 > 80 > 9000 by score, got:\n%s", out)
	}
}

func TestPrint_ShowsOverride(t *testing.T) {
	r := New()
	_ = r.SetOverride(8080, 99)
	rp := NewReporter(r)
	var buf bytes.Buffer
	rp.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "yes") {
		t.Errorf("expected 'yes' override marker, got: %s", out)
	}
}
