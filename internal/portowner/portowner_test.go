package portowner

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	r := New()
	if err := r.Set(8080, "team-backend"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	owner, ok := r.Get(8080)
	if !ok || owner != "team-backend" {
		t.Fatalf("expected team-backend, got %q ok=%v", owner, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Set(0, "team"); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
	if err := r.Set(70000, "team"); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_EmptyOwner(t *testing.T) {
	r := New()
	if err := r.Set(443, ""); err != ErrEmptyOwner {
		t.Fatalf("expected ErrEmptyOwner, got %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	r := New()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove(t *testing.T) {
	r := New()
	_ = r.Set(22, "ops")
	r.Remove(22)
	_, ok := r.Get(22)
	if ok {
		t.Fatal("expected port to be removed")
	}
}

func TestRemove_NotPresent(t *testing.T) {
	r := New()
	r.Remove(1234) // should not panic
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := New()
	_ = r.Set(80, "web")
	_ = r.Set(443, "web-tls")
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating returned map should not affect registry
	delete(all, 80)
	_, ok := r.Get(80)
	if !ok {
		t.Fatal("registry should still have port 80")
	}
}
