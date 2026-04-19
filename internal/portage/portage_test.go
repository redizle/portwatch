package portage

import (
	"testing"
	"time"
)

func TestMark_CreatesEntry(t *testing.T) {
	tr := New()
	if err := tr.Mark(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
	if e.Since.IsZero() {
		t.Error("Since should not be zero")
	}
}

func TestMark_PreservesSince(t *testing.T) {
	tr := New()
	_ = tr.Mark(443)
	e1, _ := tr.Get(443)
	time.Sleep(2 * time.Millisecond)
	_ = tr.Mark(443)
	e2, _ := tr.Get(443)
	if !e1.Since.Equal(e2.Since) {
		t.Error("Since should not change on subsequent Mark")
	}
	if !e2.LastSeen.After(e1.LastSeen) {
		t.Error("LastSeen should advance on subsequent Mark")
	}
}

func TestMark_InvalidPort(t *testing.T) {
	tr := New()
	if err := tr.Mark(0); err == nil {
		t.Error("expected error for port 0")
	}
	if err := tr.Mark(65536); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestReset_ClearsSince(t *testing.T) {
	tr := New()
	_ = tr.Mark(22)
	tr.Reset(22)
	if _, ok := tr.Get(22); ok {
		t.Error("expected entry to be removed after Reset")
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New()
	if _, ok := tr.Get(9999); ok {
		t.Error("expected no entry for untracked port")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := New()
	_ = tr.Mark(80)
	_ = tr.Mark(443)
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAge_IsPositive(t *testing.T) {
	tr := New()
	_ = tr.Mark(3000)
	time.Sleep(1 * time.Millisecond)
	e, _ := tr.Get(3000)
	if e.Age() <= 0 {
		t.Error("expected positive age")
	}
}
