package portwindow_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwindow"
)

func TestObserve_CreatesWindow(t *testing.T) {
	tr, _ := portwindow.New(time.Minute)
	now := time.Now()
	if err := tr.Observe(8080, now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected window")
	}
	if w.Hits != 1 {
		t.Errorf("expected 1 hit, got %d", w.Hits)
	}
}

func TestObserve_AccumulatesHits(t *testing.T) {
	tr, _ := portwindow.New(time.Minute)
	now := time.Now()
	_ = tr.Observe(443, now)
	_ = tr.Observe(443, now.Add(5*time.Second))
	_ = tr.Observe(443, now.Add(10*time.Second))
	w, _ := tr.Get(443)
	if w.Hits != 3 {
		t.Errorf("expected 3 hits, got %d", w.Hits)
	}
	if w.Duration != 10*time.Second {
		t.Errorf("expected 10s duration, got %v", w.Duration)
	}
}

func TestObserve_ResetsAfterWindowExpiry(t *testing.T) {
	tr, _ := portwindow.New(5 * time.Second)
	now := time.Now()
	_ = tr.Observe(22, now)
	_ = tr.Observe(22, now.Add(10*time.Second)) // beyond window
	w, _ := tr.Get(22)
	if w.Hits != 1 {
		t.Errorf("expected reset to 1 hit, got %d", w.Hits)
	}
}

func TestObserve_InvalidPort(t *testing.T) {
	tr, _ := portwindow.New(time.Minute)
	if err := tr.Observe(0, time.Now()); err == nil {
		t.Error("expected error for port 0")
	}
	if err := tr.Observe(70000, time.Now()); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	tr, _ := portwindow.New(time.Minute)
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected missing")
	}
}

func TestReset_ClearsWindow(t *testing.T) {
	tr, _ := portwindow.New(time.Minute)
	_ = tr.Observe(80, time.Now())
	tr.Reset(80)
	if tr.Len() != 0 {
		t.Error("expected empty tracker after reset")
	}
}

func TestNew_InvalidSize(t *testing.T) {
	_, err := portwindow.New(0)
	if err == nil {
		t.Error("expected error for zero window size")
	}
}
