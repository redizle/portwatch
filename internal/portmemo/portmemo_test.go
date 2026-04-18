package portmemo_test

import (
	"testing"
	"time"

	"portwatch/internal/portmemo"
)

func TestSet_And_Get(t *testing.T) {
	s := portmemo.New()
	if err := s.Set(8080, "testing http", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected memo to exist")
	}
	if m.Text != "testing http" {
		t.Errorf("got %q, want %q", m.Text, "testing http")
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := portmemo.New()
	if err := s.Set(0, "bad", 0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Set(70000, "bad", 0); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestSet_EmptyText(t *testing.T) {
	s := portmemo.New()
	if err := s.Set(443, "", 0); err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestGet_Missing(t *testing.T) {
	s := portmemo.New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected miss")
	}
}

func TestRemove(t *testing.T) {
	s := portmemo.New()
	_ = s.Set(22, "ssh note", 0)
	s.Remove(22)
	_, ok := s.Get(22)
	if ok {
		t.Fatal("expected memo to be removed")
	}
}

func TestMemo_Expires(t *testing.T) {
	s := portmemo.New()
	_ = s.Set(3306, "mysql", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, ok := s.Get(3306)
	if ok {
		t.Fatal("expected memo to be expired")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	s := portmemo.New()
	_ = s.Set(111, "expires", 10*time.Millisecond)
	_ = s.Set(222, "forever", 0)
	time.Sleep(20 * time.Millisecond)
	n := s.Purge()
	if n != 1 {
		t.Errorf("expected 1 purged, got %d", n)
	}
	_, ok := s.Get(222)
	if !ok {
		t.Fatal("non-expired memo should still exist")
	}
}
