package portquota

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	q := New()
	if err := q.Set(8080, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := q.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Limit != 10 {
		t.Errorf("expected limit 10, got %d", e.Limit)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	q := New()
	if err := q.Set(0, 5); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_InvalidQuota(t *testing.T) {
	q := New()
	if err := q.Set(80, 0); err != ErrInvalidQuota {
		t.Errorf("expected ErrInvalidQuota, got %v", err)
	}
}

func TestInc_IncrementsHits(t *testing.T) {
	q := New()
	_ = q.Set(443, 3)
	_ = q.Inc(443)
	_ = q.Inc(443)
	e, _ := q.Get(443)
	if e.Hits != 2 {
		t.Errorf("expected 2 hits, got %d", e.Hits)
	}
}

func TestInc_InvalidPort(t *testing.T) {
	q := New()
	if err := q.Inc(99999); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestExceeded_True(t *testing.T) {
	q := New()
	_ = q.Set(22, 2)
	_ = q.Inc(22)
	_ = q.Inc(22)
	e, _ := q.Get(22)
	if !e.Exceeded() {
		t.Error("expected quota to be exceeded")
	}
}

func TestExceeded_False(t *testing.T) {
	q := New()
	_ = q.Set(22, 5)
	_ = q.Inc(22)
	e, _ := q.Get(22)
	if e.Exceeded() {
		t.Error("expected quota not exceeded")
	}
}

func TestReset_ClearsHits(t *testing.T) {
	q := New()
	_ = q.Set(3000, 10)
	_ = q.Inc(3000)
	q.Reset(3000)
	e, _ := q.Get(3000)
	if e.Hits != 0 {
		t.Errorf("expected 0 hits after reset, got %d", e.Hits)
	}
	if e.Limit != 10 {
		t.Errorf("expected limit preserved as 10, got %d", e.Limit)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	q := New()
	_ = q.Set(9090, 5)
	q.Remove(9090)
	_, ok := q.Get(9090)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestGet_Missing(t *testing.T) {
	q := New()
	_, ok := q.Get(1234)
	if ok {
		t.Error("expected missing entry")
	}
}
