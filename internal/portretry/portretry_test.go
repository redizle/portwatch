package portretry

import (
	"testing"
	"time"
)

func TestRecord_IncrementsAttempts(t *testing.T) {
	r := New(100*time.Millisecond, 3)
	if err := r.Record(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := r.Get(8080)
	if e == nil || e.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %v", e)
	}
}

func TestRecord_InvalidPort(t *testing.T) {
	r := New(time.Second, 3)
	if err := r.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Record(70000); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestRecord_Accumulates(t *testing.T) {
	r := New(50*time.Millisecond, 5)
	for i := 0; i < 3; i++ {
		_ = r.Record(443)
	}
	e := r.Get(443)
	if e.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", e.Attempts)
	}
}

func TestRecord_SetsNextRetry(t *testing.T) {
	r := New(200*time.Millisecond, 3)
	before := time.Now()
	_ = r.Record(22)
	e := r.Get(22)
	if !e.NextRetry.After(before) {
		t.Fatal("expected NextRetry to be in the future")
	}
}

func TestGet_Missing(t *testing.T) {
	r := New(time.Second, 3)
	if r.Get(9999) != nil {
		t.Fatal("expected nil for unknown port")
	}
}

func TestExceeded_False(t *testing.T) {
	r := New(time.Millisecond, 3)
	_ = r.Record(80)
	if r.Exceeded(80) {
		t.Fatal("should not be exceeded after 1 attempt")
	}
}

func TestExceeded_True(t *testing.T) {
	r := New(time.Millisecond, 2)
	_ = r.Record(80)
	_ = r.Record(80)
	if !r.Exceeded(80) {
		t.Fatal("should be exceeded after 2 attempts with max=2")
	}
}

func TestReset_ClearsState(t *testing.T) {
	r := New(time.Millisecond, 3)
	_ = r.Record(3000)
	r.Reset(3000)
	if r.Get(3000) != nil {
		t.Fatal("expected nil after reset")
	}
	if r.Exceeded(3000) {
		t.Fatal("should not be exceeded after reset")
	}
}
