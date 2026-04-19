package portchain

import (
	"errors"
	"testing"
)

func TestUse_NilHandler(t *testing.T) {
	c := New()
	if err := c.Use(nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestUse_AddsHandler(t *testing.T) {
	c := New()
	_ = c.Use(func(port int, status string) error { return nil })
	if c.Len() != 1 {
		t.Fatalf("expected 1 handler, got %d", c.Len())
	}
}

func TestRun_InvalidPort(t *testing.T) {
	c := New()
	if err := c.Run(0, "open"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Run(70000, "open"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestRun_NoHandlers(t *testing.T) {
	c := New()
	if err := c.Run(8080, "open"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_AllHandlersCalled(t *testing.T) {
	c := New()
	called := 0
	for i := 0; i < 3; i++ {
		_ = c.Use(func(port int, status string) error {
			called++
			return nil
		})
	}
	_ = c.Run(443, "open")
	if called != 3 {
		t.Fatalf("expected 3 calls, got %d", called)
	}
}

func TestRun_StopsOnError(t *testing.T) {
	c := New()
	called := 0
	_ = c.Use(func(port int, status string) error { called++; return errors.New("stop") })
	_ = c.Use(func(port int, status string) error { called++; return nil })
	err := c.Run(80, "open")
	if err == nil {
		t.Fatal("expected error")
	}
	if called != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestReset_ClearsHandlers(t *testing.T) {
	c := New()
	_ = c.Use(func(port int, status string) error { return nil })
	c.Reset()
	if c.Len() != 0 {
		t.Fatalf("expected 0 handlers after reset, got %d", c.Len())
	}
}
