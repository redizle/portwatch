package portroute_test

import (
	"testing"

	"github.com/user/portwatch/internal/portroute"
)

func TestSet_And_Get(t *testing.T) {
	r := portroute.New()
	if err := r.Set(8080, "/api/v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := r.Get(8080)
	if !ok {
		t.Fatal("expected route to be found")
	}
	if got != "/api/v1" {
		t.Errorf("got %q, want /api/v1", got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := portroute.New()
	if err := r.Set(0, "/api"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Set(70000, "/api"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestSet_EmptyRoute(t *testing.T) {
	r := portroute.New()
	if err := r.Set(443, ""); err == nil {
		t.Fatal("expected error for empty route")
	}
}

func TestGet_Missing(t *testing.T) {
	r := portroute.New()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected missing for unknown port")
	}
}

func TestRemove(t *testing.T) {
	r := portroute.New()
	_ = r.Set(3000, "/frontend")
	r.Remove(3000)
	_, ok := r.Get(3000)
	if ok {
		t.Fatal("expected route to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := portroute.New()
	_ = r.Set(80, "/http")
	_ = r.Set(443, "/https")
	entries := r.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLen(t *testing.T) {
	r := portroute.New()
	if r.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	_ = r.Set(22, "/ssh")
	if r.Len() != 1 {
		t.Fatalf("expected 1, got %d", r.Len())
	}
}
