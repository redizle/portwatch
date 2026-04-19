package portcategory

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(80, "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cat, ok := s.Get(80)
	if !ok || cat != "web" {
		t.Fatalf("expected 'web', got %q ok=%v", cat, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "web"); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
	if err := s.Set(70000, "web"); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_EmptyCategory(t *testing.T) {
	s := New()
	if err := s.Set(80, ""); err != ErrEmptyCategory {
		t.Fatalf("expected ErrEmptyCategory, got %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(443)
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(22, "infra")
	s.Remove(22)
	_, ok := s.Get(22)
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, "web")
	_ = s.Set(22, "infra")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating copy should not affect store
	delete(all, 80)
	if _, ok := s.Get(80); !ok {
		t.Fatal("store should not be affected by mutation of copy")
	}
}

func TestByCategory_ReturnsMatchingPorts(t *testing.T) {
	s := New()
	_ = s.Set(80, "web")
	_ = s.Set(443, "web")
	_ = s.Set(22, "infra")
	ports := s.ByCategory("web")
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports for 'web', got %d", len(ports))
	}
}

func TestByCategory_Missing(t *testing.T) {
	s := New()
	ports := s.ByCategory("unknown")
	if len(ports) != 0 {
		t.Fatalf("expected empty slice, got %v", ports)
	}
}
