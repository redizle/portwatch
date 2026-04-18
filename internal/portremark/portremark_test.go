package portremark_test

import (
	"testing"
	"time"

	"portwatch/internal/portremark"
)

func TestSet_And_Get(t *testing.T) {
	s := portremark.New()
	if err := s.Set(8080, "needs review"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected remark to exist")
	}
	if r.Text != "needs review" {
		t.Errorf("got %q, want %q", r.Text, "needs review")
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := portremark.New()
	if err := s.Set(0, "bad"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := s.Set(99999, "bad"); err == nil {
		t.Error("expected error for port 99999")
	}
}

func TestSet_EmptyText(t *testing.T) {
	s := portremark.New()
	if err := s.Set(443, ""); err == nil {
		t.Error("expected error for empty text")
	}
}

func TestGet_Missing(t *testing.T) {
	s := portremark.New()
	_, ok := s.Get(1234)
	if ok {
		t.Error("expected missing remark")
	}
}

func TestRemove(t *testing.T) {
	s := portremark.New()
	_ = s.Set(22, "ssh remark")
	s.Remove(22)
	_, ok := s.Get(22)
	if ok {
		t.Error("expected remark to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := portremark.New()
	_ = s.Set(80, "http")
	_ = s.Set(443, "https")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("got %d entries, want 2", len(all))
	}
	// mutating copy must not affect store
	delete(all, 80)
	if _, ok := s.Get(80); !ok {
		t.Error("store should not be affected by copy mutation")
	}
}

func TestSet_OverwritesExisting(t *testing.T) {
	s := portremark.New()
	_ = s.Set(3000, "first")
	time.Sleep(time.Millisecond)
	_ = s.Set(3000, "second")
	r, _ := s.Get(3000)
	if r.Text != "second" {
		t.Errorf("expected updated text, got %q", r.Text)
	}
}
