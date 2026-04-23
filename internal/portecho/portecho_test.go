package portecho

import (
	"testing"
	"time"
)

func TestRecord_And_Get(t *testing.T) {
	e := New()
	if err := e.Record(8080, 5*time.Millisecond, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := e.Get(8080)
	if !ok {
		t.Fatal("expected result to be present")
	}
	if r.Port != 8080 {
		t.Errorf("expected port 8080, got %d", r.Port)
	}
	if r.Latency != 5*time.Millisecond {
		t.Errorf("expected latency 5ms, got %v", r.Latency)
	}
	if !r.Responded {
		t.Error("expected Responded=true")
	}
}

func TestRecord_InvalidPort(t *testing.T) {
	e := New()
	for _, p := range []int{0, -1, 65536, 99999} {
		if err := e.Record(p, 0, false); err == nil {
			t.Errorf("expected error for port %d", p)
		}
	}
}

func TestGet_Missing(t *testing.T) {
	e := New()
	_, ok := e.Get(443)
	if ok {
		t.Error("expected missing result")
	}
}

func TestRecord_Overwrites(t *testing.T) {
	e := New()
	_ = e.Record(22, 10*time.Millisecond, true)
	_ = e.Record(22, 2*time.Millisecond, false)
	r, _ := e.Get(22)
	if r.Latency != 2*time.Millisecond {
		t.Errorf("expected updated latency, got %v", r.Latency)
	}
	if r.Responded {
		t.Error("expected Responded=false after overwrite")
	}
}

func TestClear_RemovesResult(t *testing.T) {
	e := New()
	_ = e.Record(3306, 1*time.Millisecond, true)
	e.Clear(3306)
	_, ok := e.Get(3306)
	if ok {
		t.Error("expected result to be cleared")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	e := New()
	if e.Len() != 0 {
		t.Fatalf("expected 0, got %d", e.Len())
	}
	_ = e.Record(80, 0, true)
	_ = e.Record(443, 0, true)
	if e.Len() != 2 {
		t.Errorf("expected 2, got %d", e.Len())
	}
	e.Clear(80)
	if e.Len() != 1 {
		t.Errorf("expected 1 after clear, got %d", e.Len())
	}
}

func TestResult_RecordedAt_IsSet(t *testing.T) {
	e := New()
	before := time.Now()
	_ = e.Record(9000, 3*time.Millisecond, true)
	after := time.Now()
	r, _ := e.Get(9000)
	if r.RecordedAt.Before(before) || r.RecordedAt.After(after) {
		t.Errorf("RecordedAt %v not in expected range", r.RecordedAt)
	}
}
