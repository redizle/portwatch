package portrank

import (
	"testing"
)

func TestAdd_And_Get(t *testing.T) {
	r := New()
	if err := r.Add(8080, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := r.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Score != 5 {
		t.Errorf("expected score 5, got %d", e.Score)
	}
	if e.Override {
		t.Error("expected Override=false")
	}
}

func TestAdd_Accumulates(t *testing.T) {
	r := New()
	_ = r.Add(443, 3)
	_ = r.Add(443, 7)
	e, _ := r.Get(443)
	if e.Score != 10 {
		t.Errorf("expected 10, got %d", e.Score)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Add(0, 1); err == nil {
		t.Error("expected error for port 0")
	}
	if err := r.Add(70000, 1); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestSetOverride_TakesPrecedence(t *testing.T) {
	r := New()
	_ = r.Add(80, 100)
	_ = r.SetOverride(80, 1)
	e, ok := r.Get(80)
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Score != 1 {
		t.Errorf("expected override score 1, got %d", e.Score)
	}
	if !e.Override {
		t.Error("expected Override=true")
	}
}

func TestClearOverride_FallsBackToActivity(t *testing.T) {
	r := New()
	_ = r.Add(22, 9)
	_ = r.SetOverride(22, 0)
	r.ClearOverride(22)
	e, _ := r.Get(22)
	if e.Score != 9 {
		t.Errorf("expected activity score 9, got %d", e.Score)
	}
	if e.Override {
		t.Error("expected Override=false after clear")
	}
}

func TestGet_Missing(t *testing.T) {
	r := New()
	_, ok := r.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	r := New()
	_ = r.Add(3000, 5)
	r.Reset(3000)
	_, ok := r.Get(3000)
	if ok {
		t.Error("expected entry to be cleared after Reset")
	}
}

func TestLen_CountsUniquePorts(t *testing.T) {
	r := New()
	_ = r.Add(80, 1)
	_ = r.Add(443, 1)
	_ = r.SetOverride(80, 99) // same port, should not double-count
	if r.Len() != 2 {
		t.Errorf("expected Len=2, got %d", r.Len())
	}
}
