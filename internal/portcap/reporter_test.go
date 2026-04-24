package portcap

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NoPeaks(t *testing.T) {
	c := New()
	var buf bytes.Buffer
	r := NewReporter(c, &buf)
	r.Print()
	if !strings.Contains(buf.String(), "no peak") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrint_WithPeaks(t *testing.T) {
	c := New()
	_ = c.Observe(80, 15)
	_ = c.Observe(443, 200)
	var buf bytes.Buffer
	r := NewReporter(c, &buf)
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in output")
	}
	if !strings.Contains(out, "200") {
		t.Error("expected peak count 200 in output")
	}
}

func TestPrint_SortedByPort(t *testing.T) {
	c := New()
	_ = c.Observe(9000, 1)
	_ = c.Observe(22, 5)
	_ = c.Observe(3306, 3)
	var buf bytes.Buffer
	r := NewReporter(c, &buf)
	r.Print()
	out := buf.String()
	idx22 := strings.Index(out, "22")
	idx3306 := strings.Index(out, "3306")
	idx9000 := strings.Index(out, "9000")
	if idx22 > idx3306 || idx3306 > idx9000 {
		t.Error("ports are not printed in ascending order")
	}
}

func TestPrint_ContainsHeader(t *testing.T) {
	c := New()
	_ = c.Observe(8080, 7)
	var buf bytes.Buffer
	r := NewReporter(c, &buf)
	r.Print()
	if !strings.Contains(buf.String(), "PORT") {
		t.Error("expected header row in output")
	}
}
