package portcomment

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "main API"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, ok := s.Get(8080)
	if !ok || c != "main API" {
		t.Fatalf("expected 'main API', got %q ok=%v", c, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Set(70000, "bad"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected missing")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(443, "https")
	s.Remove(443)
	_, ok := s.Get(443)
	if ok {
		t.Fatal("expected removed")
	}
}

func TestRemove_NotPresent(t *testing.T) {
	s := New()
	s.Remove(1234) // should not panic
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, "http")
	_ = s.Set(22, "ssh")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the copy should not affect store
	all[80] = "changed"
	c, _ := s.Get(80)
	if c != "http" {
		t.Fatal("store was mutated")
	}
}

func TestAll_Empty(t *testing.T) {
	s := New()
	if len(s.All()) != 0 {
		t.Fatal("expected empty")
	}
}
