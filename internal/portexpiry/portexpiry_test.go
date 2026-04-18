package portexpiry

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	r := New()
	if err := r.Set(8080, time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := r.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Set(0, time.Minute); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Set(70000, time.Minute); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	r := New()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestRemove(t *testing.T) {
	r := New()
	_ = r.Set(443, time.Minute)
	r.Remove(443)
	_, ok := r.Get(443)
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestIsExpired_NotYet(t *testing.T) {
	r := New()
	_ = r.Set(80, time.Minute)
	e, _ := r.Get(80)
	if e.IsExpired() {
		t.Fatal("entry should not be expired yet")
	}
}

func TestExpired_ReturnsElapsed(t *testing.T) {
	r := New()
	r.now = func() time.Time { return time.Now().Add(-2 * time.Second) }
	_ = r.Set(22, time.Second)
	r.now = time.Now
	expired := r.Expired()
	if len(expired) != 1 || expired[0].Port != 22 {
		t.Fatalf("expected port 22 in expired list, got %v", expired)
	}
}

func TestEvict_RemovesExpired(t *testing.T) {
	r := New()
	r.now = func() time.Time { return time.Now().Add(-2 * time.Second) }
	_ = r.Set(3306, time.Second)
	r.now = time.Now
	_ = r.Set(5432, time.Minute)
	evicted := r.Evict()
	if len(evicted) != 1 || evicted[0].Port != 3306 {
		t.Fatalf("expected only port 3306 evicted, got %v", evicted)
	}
	_, ok := r.Get(3306)
	if ok {
		t.Fatal("evicted port should be removed")
	}
	_, ok = r.Get(5432)
	if !ok {
		t.Fatal("non-expired port should remain")
	}
}
