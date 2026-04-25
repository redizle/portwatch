package portweight

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	wt := New()
	if err := wt.Set(8080, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := wt.Get(8080); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestGet_DefaultFallback(t *testing.T) {
	wt := New()
	if got := wt.Get(9999); got != defaultWeight {
		t.Errorf("expected default %d, got %d", defaultWeight, got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	wt := New()
	if err := wt.Set(0, 3); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
	if err := wt.Set(70000, 3); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_InvalidWeight(t *testing.T) {
	wt := New()
	if err := wt.Set(80, 0); err != ErrInvalidWeight {
		t.Errorf("expected ErrInvalidWeight, got %v", err)
	}
	if err := wt.Set(80, -1); err != ErrInvalidWeight {
		t.Errorf("expected ErrInvalidWeight, got %v", err)
	}
}

func TestRemove_FallsBackToDefault(t *testing.T) {
	wt := New()
	_ = wt.Set(443, 10)
	wt.Remove(443)
	if got := wt.Get(443); got != defaultWeight {
		t.Errorf("expected default after remove, got %d", got)
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	wt := New()
	_ = wt.Set(80, 2)
	_ = wt.Set(443, 7)
	entries := wt.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLen_Counts(t *testing.T) {
	wt := New()
	if wt.Len() != 0 {
		t.Error("expected 0")
	}
	_ = wt.Set(22, 1)
	if wt.Len() != 1 {
		t.Error("expected 1")
	}
}

func TestAll_Empty(t *testing.T) {
	wt := New()
	if entries := wt.All(); len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}
