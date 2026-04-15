package state

import (
	"testing"
	"time"
)

func TestUpdate_NewPort(t *testing.T) {
	s := New()
	changed := s.Update(8080, true)
	if !changed {
		t.Error("expected changed=true for new port")
	}
	ps, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected port 8080 to exist")
	}
	if !ps.Open {
		t.Error("expected port to be open")
	}
}

func TestUpdate_NoChange(t *testing.T) {
	s := New()
	s.Update(9090, false)
	changed := s.Update(9090, false)
	if changed {
		t.Error("expected changed=false when status unchanged")
	}
}

func TestUpdate_StatusChange(t *testing.T) {
	s := New()
	s.Update(3000, false)
	changed := s.Update(3000, true)
	if !changed {
		t.Error("expected changed=true when status flips")
	}
	ps, _ := s.Get(3000)
	if !ps.Open {
		t.Error("expected port to now be open")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(1234)
	if ok {
		t.Error("expected ok=false for unknown port")
	}
}

func TestUpdate_TimestampsSet(t *testing.T) {
	before := time.Now()
	s := New()
	s.Update(4000, true)
	ps, _ := s.Get(4000)
	if ps.FirstSeen.Before(before) {
		t.Error("FirstSeen should be after test start")
	}
	if ps.LastSeen.Before(before) {
		t.Error("LastSeen should be after test start")
	}
}

func TestSnapshot(t *testing.T) {
	s := New()
	s.Update(80, true)
	s.Update(443, true)
	s.Update(22, false)

	snap := s.Snapshot()
	if len(snap) != 3 {
		t.Errorf("expected 3 entries in snapshot, got %d", len(snap))
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	s := New()
	s.Update(5000, true)
	snap := s.Snapshot()
	snap[0].Open = false

	ps, _ := s.Get(5000)
	if !ps.Open {
		t.Error("snapshot mutation should not affect store")
	}
}
