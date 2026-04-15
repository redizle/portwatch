package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func TestAdd_ValidPort(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(8080, "http-alt", watchlist.PriorityNormal); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Contains(8080) {
		t.Error("expected port 8080 to be in watchlist")
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(0, "bad", watchlist.PriorityLow); err == nil {
		t.Error("expected error for port 0")
	}
	if err := w.Add(65536, "bad", watchlist.PriorityLow); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestRemove(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(9000, "test", watchlist.PriorityHigh)
	w.Remove(9000)
	if w.Contains(9000) {
		t.Error("expected port 9000 to be removed")
	}
}

func TestRemove_NotPresent(t *testing.T) {
	w := watchlist.New()
	// should not panic
	w.Remove(1234)
}

func TestGet_Found(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(443, "https", watchlist.PriorityHigh)
	e, ok := w.Get(443)
	if !ok {
		t.Fatal("expected to find port 443")
	}
	if e.Label != "https" {
		t.Errorf("expected label 'https', got %q", e.Label)
	}
	if e.Priority != watchlist.PriorityHigh {
		t.Errorf("expected PriorityHigh, got %d", e.Priority)
	}
}

func TestGet_Missing(t *testing.T) {
	w := watchlist.New()
	_, ok := w.Get(9999)
	if ok {
		t.Error("expected missing port to return false")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(80, "http", watchlist.PriorityNormal)
	_ = w.Add(22, "ssh", watchlist.PriorityHigh)
	all := w.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the slice should not affect internal state
	all[0].Label = "mutated"
	for _, e := range w.All() {
		if e.Label == "mutated" {
			t.Error("internal state was mutated through returned slice")
		}
	}
}
