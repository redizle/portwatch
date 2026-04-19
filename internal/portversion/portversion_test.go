package portversion

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "1.2.3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if v != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %s", v)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "1.0.0"); err == nil {
		t.Fatal("expected error for invalid port")
	}
	if err := s.Set(70000, "1.0.0"); err == nil {
		t.Fatal("expected error for port > 65535")
	}
}

func TestSet_EmptyVersion(t *testing.T) {
	s := New()
	if err := s.Set(443, ""); err == nil {
		t.Fatal("expected error for empty version")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(22, "openssh-9.0")
	s.Remove(22)
	_, ok := s.Get(22)
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, "nginx/1.24")
	_ = s.Set(443, "nginx/1.24")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestSet_Overwrite(t *testing.T) {
	s := New()
	_ = s.Set(3306, "mysql-8.0")
	_ = s.Set(3306, "mysql-8.1")
	v, _ := s.Get(3306)
	if v != "mysql-8.1" {
		t.Fatalf("expected mysql-8.1, got %s", v)
	}
}
