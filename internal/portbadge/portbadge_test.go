package portbadge

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(80, "✔", "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b := s.Get(80)
	if b.Label != "http" || b.Icon != "✔" {
		t.Errorf("got %+v", b)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "x", "bad"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := s.Set(70000, "x", "bad"); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestSet_EmptyLabel(t *testing.T) {
	s := New()
	if err := s.Set(443, "!", ""); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestGet_Missing_ReturnsDefault(t *testing.T) {
	s := New()
	b := s.Get(9999)
	if b.Label != "unknown" {
		t.Errorf("expected default badge, got %+v", b)
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(22, "🔒", "ssh")
	s.Remove(22)
	b := s.Get(22)
	if b.Label != "unknown" {
		t.Errorf("expected default after remove, got %+v", b)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, "✔", "http")
	_ = s.Set(443, "🔒", "https")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	// mutating copy should not affect store
	delete(all, 80)
	if len(s.All()) != 2 {
		t.Error("store was mutated via All()")
	}
}
