package portpriority

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	r := New(Normal)
	if err := r.Set(8080, High); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Get(8080); got != High {
		t.Errorf("expected High, got %v", got)
	}
}

func TestGet_DefaultFallback(t *testing.T) {
	r := New(Low)
	if got := r.Get(9999); got != Low {
		t.Errorf("expected Low fallback, got %v", got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := New(Normal)
	if err := r.Set(0, High); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
	if err := r.Set(65536, High); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_InvalidLevel(t *testing.T) {
	r := New(Normal)
	if err := r.Set(80, Level(99)); err != ErrInvalidLevel {
		t.Errorf("expected ErrInvalidLevel, got %v", err)
	}
}

func TestRemove_FallsBackToDefault(t *testing.T) {
	r := New(Normal)
	_ = r.Set(443, Critical)
	r.Remove(443)
	if got := r.Get(443); got != Normal {
		t.Errorf("expected Normal after remove, got %v", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := New(Normal)
	_ = r.Set(22, High)
	_ = r.Set(3306, Critical)
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the copy should not affect the registry
	delete(all, 22)
	if got := r.Get(22); got != High {
		t.Errorf("registry mutated by copy modification")
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct{ l Level; s string }{
		{Low, "low"}, {Normal, "normal"}, {High, "high"}, {Critical, "critical"}, {Level(0), "unknown"},
	}
	for _, c := range cases {
		if got := c.l.String(); got != c.s {
			t.Errorf("Level(%d).String() = %q, want %q", c.l, got, c.s)
		}
	}
}
