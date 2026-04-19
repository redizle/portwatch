package portcooldown_test

import (
	"testing"
	"time"

	"portwatch/internal/portcooldown"
)

func TestTrigger_SetsCooldown(t *testing.T) {
	c := portcooldown.New(100 * time.Millisecond)
	if err := c.Trigger(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.IsCooling(8080) {
		t.Fatal("expected port to be cooling")
	}
}

func TestTrigger_InvalidPort(t *testing.T) {
	c := portcooldown.New(time.Second)
	if err := c.Trigger(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Trigger(70000); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestIsCooling_False_WhenNotTriggered(t *testing.T) {
	c := portcooldown.New(time.Second)
	if c.IsCooling(9090) {
		t.Fatal("expected not cooling for untriggered port")
	}
}

func TestIsCooling_Expires(t *testing.T) {
	c := portcooldown.New(10 * time.Millisecond)
	_ = c.Trigger(443)
	time.Sleep(20 * time.Millisecond)
	if c.IsCooling(443) {
		t.Fatal("expected cooldown to have expired")
	}
}

func TestReset_ClearsCooldown(t *testing.T) {
	c := portcooldown.New(time.Minute)
	_ = c.Trigger(22)
	c.Reset(22)
	if c.IsCooling(22) {
		t.Fatal("expected cooldown to be cleared after reset")
	}
}

func TestGet_TracksTriggeredCount(t *testing.T) {
	c := portcooldown.New(time.Minute)
	_ = c.Trigger(3000)
	_ = c.TriggerFor(3000, time.Minute)
	e, ok := c.Get(3000)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Triggered != 2 {
		t.Fatalf("expected 2 triggers, got %d", e.Triggered)
	}
}

func TestGet_Missing(t *testing.T) {
	c := portcooldown.New(time.Second)
	_, ok := c.Get(1234)
	if ok {
		t.Fatal("expected missing entry")
	}
}
