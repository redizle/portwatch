package portenv

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New("dev")
	if err := s.Set(8080, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get(8080); got != "prod" {
		t.Errorf("expected prod, got %s", got)
	}
}

func TestGet_DefaultFallback(t *testing.T) {
	s := New("staging")
	if got := s.Get(9999); got != "staging" {
		t.Errorf("expected staging, got %s", got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New("dev")
	if err := s.Set(0, "prod"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := s.Set(70000, "prod"); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestSet_EmptyEnv(t *testing.T) {
	s := New("dev")
	if err := s.Set(80, ""); err == nil {
		t.Error("expected error for empty env")
	}
}

func TestRemove_FallsBackToDefault(t *testing.T) {
	s := New("dev")
	_ = s.Set(443, "prod")
	s.Remove(443)
	if got := s.Get(443); got != "dev" {
		t.Errorf("expected dev after remove, got %s", got)
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	s := New("dev")
	_ = s.Set(80, "prod")
	_ = s.Set(22, "staging")
	entries := s.All()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_Empty(t *testing.T) {
	s := New("dev")
	if got := s.All(); len(got) != 0 {
		t.Errorf("expected empty, got %d entries", len(got))
	}
}
