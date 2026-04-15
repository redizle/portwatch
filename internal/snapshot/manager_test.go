package snapshot

import (
	"testing"
)

func TestNewManager_HasEmptyCurrent(t *testing.T) {
	m := NewManager()
	curr := m.Current()
	if curr == nil {
		t.Fatal("expected non-nil current snapshot")
	}
	if len(curr.All()) != 0 {
		t.Error("expected empty initial snapshot")
	}
}

func TestRotate_PromotesCurrent(t *testing.T) {
	m := NewManager()
	m.Current().Set(8080, true)

	changed := m.Rotate()

	// first rotation from nil prev: open ports are reported as changed
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed port on first rotate, got %d", len(changed))
	}

	prev := m.Previous()
	if prev == nil {
		t.Fatal("expected previous snapshot to be set after rotate")
	}

	ps, ok := prev.Get(8080)
	if !ok || !ps.Open {
		t.Error("expected port 8080 open in previous snapshot")
	}
}

func TestRotate_FreshCurrentAfterRotate(t *testing.T) {
	m := NewManager()
	m.Current().Set(9000, true)
	m.Rotate()

	curr := m.Current()
	if len(curr.All()) != 0 {
		t.Error("expected fresh empty snapshot after rotate")
	}
}

func TestRotate_NoDiffWhenUnchanged(t *testing.T) {
	m := NewManager()
	m.Current().Set(80, true)
	m.Rotate() // first cycle establishes baseline

	m.Current().Set(80, true)
	changed := m.Rotate() // second cycle, same state

	if len(changed) != 0 {
		t.Errorf("expected no changes on second rotate, got %d", len(changed))
	}
}

func TestPrevious_NilBeforeFirstRotate(t *testing.T) {
	m := NewManager()
	if m.Previous() != nil {
		t.Error("expected nil previous before first rotate")
	}
}
