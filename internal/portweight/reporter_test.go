package portweight

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoEntries(t *testing.T) {
	wt := New()
	r := NewReporter(wt)
	var buf bytes.Buffer
	r.Print(&buf)
	if !strings.Contains(buf.String(), "no port weights") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrint_WithEntries(t *testing.T) {
	wt := New()
	_ = wt.Set(80, 3)
	_ = wt.Set(443, 9)
	r := NewReporter(wt)
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
	if !strings.Contains(out, "9") {
		t.Error("expected weight 9 in output")
	}
}

func TestPrint_SortedByPort(t *testing.T) {
	wt := New()
	_ = wt.Set(9000, 1)
	_ = wt.Set(22, 5)
	_ = wt.Set(3306, 2)
	r := NewReporter(wt)
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()
	pos22 := strings.Index(out, "22")
	pos3306 := strings.Index(out, "3306")
	pos9000 := strings.Index(out, "9000")
	if pos22 > pos3306 || pos3306 > pos9000 {
		t.Errorf("expected ports sorted ascending, got:\n%s", out)
	}
}
