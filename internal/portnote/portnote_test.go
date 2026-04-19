package portnote

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "watching closely"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected note to exist")
	}
	if n.Text != "watching closely" {
		t.Errorf("got %q", n.Text)
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

func TestSet_EmptyText(t *testing.T) {
	s := New()
	if err := s.Set(443, ""); err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected missing")
	}
}

func TestSet_UpdatesTimestamp(t *testing.T) {
	s := New()
	_ = s.Set(22, "first")
	n1, _ := s.Get(22)
	time.Sleep(2 * time.Millisecond)
	_ = s.Set(22, "second")
	n2, _ := s.Get(22)
	if !n2.UpdatedAt.After(n1.UpdatedAt) {
		t.Error("UpdatedAt should advance on update")
	}
	if n2.CreatedAt != n1.CreatedAt {
		t.Error("CreatedAt should not change on update")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(3306, "mysql")
	s.Remove(3306)
	_, ok := s.Get(3306)
	if ok {
		t.Fatal("expected note to be removed")
	}
}

func TestLen(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Fatal("expected 0")
	}
	_ = s.Set(80, "http")
	_ = s.Set(443, "https")
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}
