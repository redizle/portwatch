package portlock

import (
	"testing"
	"time"
)

func TestLock_ValidPort(t *testing.T) {
	l := New()
	if err := l.Lock(8080, "testing", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.IsLocked(8080) {
		t.Fatal("expected port 8080 to be locked")
	}
}

func TestLock_InvalidPort(t *testing.T) {
	l := New()
	if err := l.Lock(0, "bad", nil); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := l.Lock(65536, "bad", nil); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestUnlock_RemovesLock(t *testing.T) {
	l := New()
	_ = l.Lock(443, "ssl", nil)
	l.Unlock(443)
	if l.IsLocked(443) {
		t.Fatal("expected port 443 to be unlocked")
	}
}

func TestIsLocked_MissingPort(t *testing.T) {
	l := New()
	if l.IsLocked(9999) {
		t.Fatal("expected false for unknown port")
	}
}

func TestLock_WithTTL_Expires(t *testing.T) {
	l := New()
	ttl := 30 * time.Millisecond
	_ = l.Lock(3000, "temp", &ttl)
	if !l.IsLocked(3000) {
		t.Fatal("expected port to be locked before expiry")
	}
	time.Sleep(50 * time.Millisecond)
	if l.IsLocked(3000) {
		t.Fatal("expected port to be expired")
	}
}

func TestActive_ReturnsLiveLocks(t *testing.T) {
	l := New()
	_ = l.Lock(80, "http", nil)
	_ = l.Lock(22, "ssh", nil)
	active := l.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active locks, got %d", len(active))
	}
}

func TestActive_ExcludesExpired(t *testing.T) {
	l := New()
	ttl := 20 * time.Millisecond
	_ = l.Lock(5000, "expire-me", &ttl)
	_ = l.Lock(5001, "keep", nil)
	time.Sleep(40 * time.Millisecond)
	active := l.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active lock, got %d", len(active))
	}
	if active[0].Port != 5001 {
		t.Fatalf("expected port 5001, got %d", active[0].Port)
	}
}
