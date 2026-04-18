package portpause

import (
	"testing"
	"time"
)

func TestPause_And_IsPaused(t *testing.T) {
	p := New()
	if err := p.Pause(8080, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.IsPaused(8080) {
		t.Fatal("expected port to be paused")
	}
}

func TestPause_InvalidPort(t *testing.T) {
	p := New()
	if err := p.Pause(0, 0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := p.Pause(70000, 0); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestResume_RemovesPause(t *testing.T) {
	p := New()
	_ = p.Pause(9090, 0)
	p.Resume(9090)
	if p.IsPaused(9090) {
		t.Fatal("expected port to be resumed")
	}
}

func TestIsPaused_Missing(t *testing.T) {
	p := New()
	if p.IsPaused(1234) {
		t.Fatal("expected false for unknown port")
	}
}

func TestPause_Expires(t *testing.T) {
	p := New()
	_ = p.Pause(3000, 20*time.Millisecond)
	if !p.IsPaused(3000) {
		t.Fatal("expected port to be paused immediately")
	}
	time.Sleep(40 * time.Millisecond)
	if p.IsPaused(3000) {
		t.Fatal("expected pause to have expired")
	}
}

func TestList_ReturnsPausedPorts(t *testing.T) {
	p := New()
	_ = p.Pause(80, 0)
	_ = p.Pause(443, time.Minute)
	ports := p.List()
	if len(ports) != 2 {
		t.Fatalf("expected 2 paused ports, got %d", len(ports))
	}
}

func TestList_ExcludesExpired(t *testing.T) {
	p := New()
	_ = p.Pause(8080, 10*time.Millisecond)
	time.Sleep(30 * time.Millisecond)
	// trigger expiry via IsPaused
	p.IsPaused(8080)
	if len(p.List()) != 0 {
		t.Fatal("expected empty list after expiry")
	}
}
