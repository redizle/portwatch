package trend

import (
	"testing"
	"time"
)

func TestChurn_NoEvents(t *testing.T) {
	tr := New(time.Minute, 50)
	if c := tr.Churn(8080); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestChurn_CountsWithinWindow(t *testing.T) {
	tr := New(time.Minute, 50)
	tr.Record(8080, true)
	tr.Record(8080, false)
	tr.Record(8080, true)
	if c := tr.Churn(8080); c != 3 {
		t.Fatalf("expected 3, got %d", c)
	}
}

func TestChurn_ExcludesOutsideWindow(t *testing.T) {
	tr := New(50*time.Millisecond, 50)
	tr.Record(9000, true)
	time.Sleep(80 * time.Millisecond)
	tr.Record(9000, false)
	// only the second event should be within window
	if c := tr.Churn(9000); c != 1 {
		t.Fatalf("expected 1, got %d", c)
	}
}

func TestFlapping_BelowThreshold(t *testing.T) {
	tr := New(time.Minute, 50)
	tr.Record(443, true)
	tr.Record(443, false)
	if tr.Flapping(443, 5) {
		t.Fatal("expected not flapping")
	}
}

func TestFlapping_AtThreshold(t *testing.T) {
	tr := New(time.Minute, 50)
	for i := 0; i < 5; i++ {
		tr.Record(443, i%2 == 0)
	}
	if !tr.Flapping(443, 5) {
		t.Fatal("expected flapping")
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	tr := New(time.Minute, 50)
	tr.Record(80, true)
	tr.Record(80, false)
	tr.Reset(80)
	if c := tr.Churn(80); c != 0 {
		t.Fatalf("expected 0 after reset, got %d", c)
	}
}

func TestMaxPer_CapsStoredEvents(t *testing.T) {
	tr := New(time.Minute, 3)
	for i := 0; i < 10; i++ {
		tr.Record(22, i%2 == 0)
	}
	tr.mu.Lock()
	l := len(tr.events[22])
	tr.mu.Unlock()
	if l != 3 {
		t.Fatalf("expected 3 stored events, got %d", l)
	}
}
