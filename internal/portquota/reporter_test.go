package portquota

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoEntries(t *testing.T) {
	q := New()
	var buf bytes.Buffer
	NewReporter(q, &buf).Print()
	if !strings.Contains(buf.String(), "no quota entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrint_WithEntries(t *testing.T) {
	q := New()
	_ = q.Set(80, 5)
	_ = q.Inc(80)
	_ = q.Inc(80)
	_ = q.Set(443, 3)
	_ = q.Inc(443)
	_ = q.Inc(443)
	_ = q.Inc(443)

	var buf bytes.Buffer
	NewReporter(q, &buf).Print()
	out := buf.String()

	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
	if !strings.Contains(out, "YES") {
		t.Error("expected exceeded marker for port 443")
	}
}

func TestPrint_SortedByPort(t *testing.T) {
	q := New()
	_ = q.Set(9000, 1)
	_ = q.Set(80, 1)
	_ = q.Set(3000, 1)

	var buf bytes.Buffer
	NewReporter(q, &buf).Print()
	out := buf.String()

	i80 := strings.Index(out, "80")
	i3000 := strings.Index(out, "3000")
	i9000 := strings.Index(out, "9000")
	if !(i80 < i3000 && i3000 < i9000) {
		t.Error("expected output sorted by port number")
	}
}
