package portwatch

import (
	"testing"
	"time"
)

func TestTouch_CreatesEntry(t *testing.T) {
	w := New()
	w.Touch(8080, true)
	e, ok := w.Get(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if !e.Open {
		t.Error("expected Open=true")
	}
	if e.SeenCount != 1 {
		t.Errorf("expected SeenCount=1, got %d", e.SeenCount)
	}
}

func TestTouch_Accumulates(t *testing.T) {
	w := New()
	w.Touch(9000, true)
	w.Touch(9000, false)
	e, _ := w.Get(9000)
	if e.SeenCount != 2 {
		t.Errorf("expected SeenCount=2, got %d", e.SeenCount)
	}
	if e.Open {
		t.Error("expected Open=false after second touch")
	}
}

func TestTouch_SetsTimestamps(t *testing.T) {
	before := time.Now()
	w := New()
	w.Touch(443, true)
	after := time.Now()
	e, _ := w.Get(443)
	if e.FirstSeen.Before(before) || e.FirstSeen.After(after) {
		t.Error("FirstSeen out of expected range")
	}
}

func TestSetLabel_Valid(t *testing.T) {
	w := New()
	if err := w.SetLabel(80, "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := w.Get(80)
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Label != "http" {
		t.Errorf("expected label 'http', got %q", e.Label)
	}
}

func TestSetLabel_InvalidPort(t *testing.T) {
	w := New()
	if err := w.SetLabel(0, "bad"); err == nil {
		t.Error("expected error for port 0")
	}
}

func TestSetOwner_Valid(t *testing.T) {
	w := New()
	_ = w.SetOwner(22, "alice")
	e, _ := w.Get(22)
	if e.Owner != "alice" {
		t.Errorf("expected owner 'alice', got %q", e.Owner)
	}
}

func TestGet_Missing(t *testing.T) {
	w := New()
	_, ok := w.Get(1234)
	if ok {
		t.Error("expected missing entry")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	w := New()
	w.Touch(100, true)
	w.Touch(200, false)
	all := w.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
