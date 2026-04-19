package portcount

import (
	"testing"
)

func TestInc_And_Get(t *testing.T) {
	c := New()
	if err := c.Inc(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Get(8080); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestInc_Accumulates(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		_ = c.Inc(443)
	}
	if got := c.Get(443); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestInc_InvalidPort(t *testing.T) {
	c := New()
	if err := c.Inc(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Inc(70000); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	c := New()
	if got := c.Get(9999); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	c := New()
	_ = c.Inc(22)
	_ = c.Inc(22)
	c.Reset(22)
	if got := c.Get(22); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	c := New()
	_ = c.Inc(80)
	_ = c.Inc(443)
	all := c.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the copy should not affect internal state
	all[80] = 999
	if c.Get(80) != 1 {
		t.Fatal("internal state was mutated")
	}
}
