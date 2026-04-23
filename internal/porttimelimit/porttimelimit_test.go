package porttimelimit

import (
	"testing"
	"time"
)

func todayAt(h, m int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())
}

func TestSet_And_Get(t *testing.T) {
	l := New()
	w := Window{Start: 9 * time.Hour, End: 17 * time.Hour}
	if err := l.Set(8080, w); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := l.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got != w {
		t.Errorf("got %v, want %v", got, w)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	l := New()
	if err := l.Set(0, Window{Start: time.Hour, End: 2 * time.Hour}); err == nil {
		t.Error("expected error for port 0")
	}
}

func TestSet_InvalidWindow(t *testing.T) {
	l := New()
	if err := l.Set(80, Window{Start: 10 * time.Hour, End: 9 * time.Hour}); err == nil {
		t.Error("expected error when end <= start")
	}
}

func TestGet_Missing(t *testing.T) {
	l := New()
	_, ok := l.Get(9999)
	if ok {
		t.Error("expected missing entry")
	}
}

func TestRemove(t *testing.T) {
	l := New()
	_ = l.Set(443, Window{Start: time.Hour, End: 2 * time.Hour})
	l.Remove(443)
	_, ok := l.Get(443)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestAllowed_WithinWindow(t *testing.T) {
	l := New()
	_ = l.Set(8080, Window{Start: 9 * time.Hour, End: 17 * time.Hour})
	if !l.Allowed(8080, todayAt(12, 0)) {
		t.Error("expected port to be allowed at noon")
	}
}

func TestAllowed_OutsideWindow(t *testing.T) {
	l := New()
	_ = l.Set(8080, Window{Start: 9 * time.Hour, End: 17 * time.Hour})
	if l.Allowed(8080, todayAt(20, 0)) {
		t.Error("expected port to be blocked at 20:00")
	}
}

func TestAllowed_NoWindow_AlwaysTrue(t *testing.T) {
	l := New()
	if !l.Allowed(3000, todayAt(3, 0)) {
		t.Error("expected unregistered port to always be allowed")
	}
}

func TestViolations_ReturnsBlockedPorts(t *testing.T) {
	l := New()
	_ = l.Set(8080, Window{Start: 9 * time.Hour, End: 17 * time.Hour})
	_ = l.Set(443, Window{Start: 8 * time.Hour, End: 18 * time.Hour})
	active := []int{8080, 443, 22}
	v := l.Violations(active, todayAt(20, 0))
	if len(v) != 2 {
		t.Errorf("expected 2 violations, got %d", len(v))
	}
}

func TestViolations_NoneWhenAllAllowed(t *testing.T) {
	l := New()
	_ = l.Set(8080, Window{Start: 9 * time.Hour, End: 17 * time.Hour})
	v := l.Violations([]int{8080}, todayAt(10, 0))
	if len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}
