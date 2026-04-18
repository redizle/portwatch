package portannot

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(8080, "env")
	if !ok || v != "prod" {
		t.Fatalf("expected prod, got %q ok=%v", v, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "k", "v"); err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestSet_EmptyKey(t *testing.T) {
	s := New()
	if err := s.Set(80, "", "v"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999, "missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(443, "tier", "web")
	s.Remove(443, "tier")
	_, ok := s.Get(443, "tier")
	if ok {
		t.Fatal("expected key to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(22, "a", "1")
	_ = s.Set(22, "b", "2")
	all := s.All(22)
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	all["a"] = "mutated"
	v, _ := s.Get(22, "a")
	if v == "mutated" {
		t.Fatal("All should return a copy")
	}
}

func TestClear_RemovesAll(t *testing.T) {
	s := New()
	_ = s.Set(3306, "db", "mysql")
	s.Clear(3306)
	if len(s.All(3306)) != 0 {
		t.Fatal("expected empty after clear")
	}
}
