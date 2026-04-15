package snapshot

import (
	"testing"
	"time"
)

func TestSet_NewPort(t *testing.T) {
	s := New()
	s.Set(8080, true)

	ps, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected port 8080 to exist")
	}
	if !ps.Open {
		t.Error("expected port to be open")
	}
	if ps.FirstSeen.IsZero() {
		t.Error("expected FirstSeen to be set")
	}
}

func TestSet_UpdatesExisting(t *testing.T) {
	s := New()
	s.Set(9090, true)

	first, _ := s.Get(9090)
	time.Sleep(2 * time.Millisecond)
	s.Set(9090, false)

	updated, _ := s.Get(9090)
	if updated.Open {
		t.Error("expected port to be closed after update")
	}
	if !updated.LastSeen.After(first.LastSeen) {
		t.Error("expected LastSeen to advance on update")
	}
	if updated.FirstSeen != first.FirstSeen {
		t.Error("expected FirstSeen to remain unchanged")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(1234)
	if ok {
		t.Error("expected missing port to return false")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	s.Set(80, true)
	s.Set(443, true)
	s.Set(22, false)

	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
}

func TestDiff_NilPrev_ReturnsOpenPorts(t *testing.T) {
	s := New()
	s.Set(80, true)
	s.Set(81, false)

	changed := s.Diff(nil)
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed port, got %d", len(changed))
	}
	if changed[0].Port != 80 {
		t.Errorf("expected port 80, got %d", changed[0].Port)
	}
}

func TestDiff_DetectsStatusChange(t *testing.T) {
	prev := New()
	prev.Set(8080, false)

	curr := New()
	curr.Set(8080, true)

	changed := curr.Diff(prev)
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed port, got %d", len(changed))
	}
}

func TestDiff_NoChange(t *testing.T) {
	prev := New()
	prev.Set(443, true)

	curr := New()
	curr.Set(443, true)

	changed := curr.Diff(prev)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %d", len(changed))
	}
}
