package portchain

import (
	"testing"
)

func TestBuilder_Build_EmptyChain(t *testing.T) {
	c, err := NewBuilder().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 handlers")
	}
}

func TestBuilder_Add_ValidHandlers(t *testing.T) {
	c, err := NewBuilder().
		Add(func(port int, status string) error { return nil }).
		Add(func(port int, status string) error { return nil }).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 2 {
		t.Fatalf("expected 2 handlers, got %d", c.Len())
	}
}

func TestBuilder_Add_NilHandler_ReturnsError(t *testing.T) {
	_, err := NewBuilder().Add(nil).Build()
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestBuilder_StopsAfterFirstError(t *testing.T) {
	good := func(port int, status string) error { return nil }
	b := NewBuilder().Add(nil).Add(good)
	_, err := b.Build()
	if err == nil {
		t.Fatal("expected retained error")
	}
}

func TestBuilder_Run_ExecutesChain(t *testing.T) {
	called := false
	c, _ := NewBuilder().
		Add(func(port int, status string) error {
			called = true
			return nil
		}).Build()
	if err := c.Run(8080, "open"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("handler was not called")
	}
}
