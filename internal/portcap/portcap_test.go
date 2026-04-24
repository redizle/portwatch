package portcap

import (
	"testing"
	"time"
)

func TestObserve_SetsPeak(t *testing.T) {
	c := New()
	if err := c.Observe(8080, 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := c.Get(8080)
	if !ok {
		t.Fatal("expected peak to exist")
	}
	if p.Count != 42 {
		t.Errorf("expected 42, got %d", p.Count)
	}
}

func TestObserve_UpdatesWhenHigher(t *testing.T) {
	c := New()
	_ = c.Observe(9000, 10)
	_ = c.Observe(9000, 50)
	p, _ := c.Get(9000)
	if p.Count != 50 {
		t.Errorf("expected 50, got %d", p.Count)
	}
}

func TestObserve_IgnoresLowerValue(t *testing.T) {
	c := New()
	_ = c.Observe(443, 100)
	_ = c.Observe(443, 5)
	p, _ := c.Get(443)
	if p.Count != 100 {
		t.Errorf("expected 100, got %d", p.Count)
	}
}

func TestObserve_InvalidPort(t *testing.T) {
	c := New()
	if err := c.Observe(0, 1); err == nil {
		t.Error("expected error for port 0")
	}
	if err := c.Observe(70000, 1); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	c := New()
	_, ok := c.Get(1234)
	if ok {
		t.Error("expected missing for unseen port")
	}
}

func TestReset_ClearsPeak(t *testing.T) {
	c := New()
	_ = c.Observe(22, 7)
	if err := c.Reset(22); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := c.Get(22)
	if ok {
		t.Error("expected peak to be cleared")
	}
}

func TestReset_InvalidPort(t *testing.T) {
	c := New()
	if err := c.Reset(-1); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	c := New()
	_ = c.Observe(80, 3)
	_ = c.Observe(443, 9)
	all := c.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	// Mutating the copy must not affect internal state.
	all[80] = Peak{Count: 999, RecordedAt: time.Now()}
	p, _ := c.Get(80)
	if p.Count == 999 {
		t.Error("copy mutation affected internal state")
	}
}

func TestObserve_RecordedAtIsSet(t *testing.T) {
	c := New()
	before := time.Now()
	_ = c.Observe(3000, 1)
	after := time.Now()
	p, _ := c.Get(3000)
	if p.RecordedAt.Before(before) || p.RecordedAt.After(after) {
		t.Error("RecordedAt is outside expected range")
	}
}
