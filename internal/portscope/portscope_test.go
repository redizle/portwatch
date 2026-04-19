package portscope

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	c := New(ScopeUnknown)
	if err := c.Set(8080, ScopeExternal); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Get(8080); got != ScopeExternal {
		t.Errorf("expected %s, got %s", ScopeExternal, got)
	}
}

func TestGet_DefaultFallback(t *testing.T) {
	c := New(ScopeInternal)
	if got := c.Get(9999); got != ScopeInternal {
		t.Errorf("expected default %s, got %s", ScopeInternal, got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	c := New(ScopeUnknown)
	if err := c.Set(0, ScopeLoopback); err == nil {
		t.Error("expected error for port 0")
	}
	if err := c.Set(70000, ScopeLoopback); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestSet_EmptyScope(t *testing.T) {
	c := New(ScopeUnknown)
	if err := c.Set(80, ""); err == nil {
		t.Error("expected error for empty scope")
	}
}

func TestRemove_FallsBackToDefault(t *testing.T) {
	c := New(ScopeLoopback)
	_ = c.Set(443, ScopeExternal)
	c.Remove(443)
	if got := c.Get(443); got != ScopeLoopback {
		t.Errorf("expected default after remove, got %s", got)
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	c := New(ScopeUnknown)
	_ = c.Set(22, ScopeInternal)
	_ = c.Set(80, ScopeExternal)
	entries := c.All()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_Empty(t *testing.T) {
	c := New(ScopeUnknown)
	if got := c.All(); len(got) != 0 {
		t.Errorf("expected empty, got %d entries", len(got))
	}
}
