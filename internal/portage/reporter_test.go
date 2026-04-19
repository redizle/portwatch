package portage

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoPorts(t *testing.T) {
	tr := New()
	var buf bytes.Buffer
	rep := NewReporter(tr, &buf)
	rep.Print()
	if !strings.Contains(buf.String(), "PORT") {
		t.Error("expected header row")
	}
}

func TestPrint_WithPorts(t *testing.T) {
	tr := New()
	_ = tr.Mark(80)
	_ = tr.Mark(443)
	var buf bytes.Buffer
	rep := NewReporter(tr, &buf)
	rep.Print()
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
}

func TestPrint_SortedByPort(t *testing.T) {
	tr := New()
	_ = tr.Mark(9000)
	_ = tr.Mark(22)
	_ = tr.Mark(8080)
	var buf bytes.Buffer
	rep := NewReporter(tr, &buf)
	rep.Print()
	out := buf.String()
	i22 := strings.Index(out, "22")
	i8080 := strings.Index(out, "8080")
	i9000 := strings.Index(out, "9000")
	if !(i22 < i8080 && i8080 < i9000) {
		t.Error("expected ports sorted ascending")
	}
}
