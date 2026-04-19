package portdepend

import (
	"testing"
)

func TestAdd_And_DepsOf(t *testing.T) {
	tr := New()
	if err := tr.Add(8080, 5432); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps := tr.DepsOf(8080)
	if len(deps) != 1 || deps[0] != 5432 {
		t.Fatalf("expected [5432], got %v", deps)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	tr := New()
	if err := tr.Add(0, 80); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
	if err := tr.Add(80, 99999); err != ErrInvalidPort {
		t.Fatalf("expected ErrInvalidPort, got %v", err)
	}
}

func TestAdd_SelfDependency(t *testing.T) {
	tr := New()
	if err := tr.Add(80, 80); err != ErrSelfDependency {
		t.Fatalf("expected ErrSelfDependency, got %v", err)
	}
}

func TestRemove_DeletesEdge(t *testing.T) {
	tr := New()
	_ = tr.Add(8080, 5432)
	tr.Remove(8080, 5432)
	if len(tr.DepsOf(8080)) != 0 {
		t.Fatal("expected no deps after remove")
	}
}

func TestDependents_ReturnsCallers(t *testing.T) {
	tr := New()
	_ = tr.Add(8080, 5432)
	_ = tr.Add(9090, 5432)
	deps := tr.Dependents(5432)
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependents, got %d", len(deps))
	}
}

func TestDependents_NoneFound(t *testing.T) {
	tr := New()
	if len(tr.Dependents(3306)) != 0 {
		t.Fatal("expected empty dependents")
	}
}

func TestClear_RemovesAll(t *testing.T) {
	tr := New()
	_ = tr.Add(8080, 5432)
	_ = tr.Add(8080, 6379)
	tr.Clear(8080)
	if len(tr.DepsOf(8080)) != 0 {
		t.Fatal("expected no deps after clear")
	}
}

func TestString_WithDeps(t *testing.T) {
	tr := New()
	_ = tr.Add(8080, 5432)
	s := tr.String(8080)
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestString_NoDeps(t *testing.T) {
	tr := New()
	s := tr.String(8080)
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
