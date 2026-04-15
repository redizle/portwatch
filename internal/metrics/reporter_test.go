package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_ContainsAllFields(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)

	c := Counters{
		ScansTotal:    42,
		PortsOpen:     7,
		PortsClosed:   3,
		AlertsSent:    1,
		FilterDropped: 5,
		StartedAt:     time.Now().Add(-10 * time.Second),
	}

	if err := r.Print(c); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	out := buf.String()
	checks := []string{
		"Uptime",
		"Scans total",
		"42",
		"Ports open",
		"7",
		"Ports closed",
		"3",
		"Alerts sent",
		"1",
		"Filter dropped",
		"5",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrint_ZeroCounters(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	c := Counters{StartedAt: time.Now()}

	if err := r.Print(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output even for zero counters")
	}
}
