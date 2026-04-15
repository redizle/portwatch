package notify

import (
	"errors"
	"testing"
	"time"
)

func TestDispatch_NoHandlers(t *testing.T) {
	n := New(0)
	if err := n.Dispatch(8080, "open", LevelInfo); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDispatch_CallsHandler(t *testing.T) {
	n := New(0)
	var received Event
	n.Register(func(e Event) error {
		received = e
		return nil
	})
	if err := n.Dispatch(9090, "open", LevelWarn); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Port != 9090 {
		t.Errorf("expected port 9090, got %d", received.Port)
	}
	if received.Status != "open" {
		t.Errorf("expected status open, got %s", received.Status)
	}
	if received.Level != LevelWarn {
		t.Errorf("expected level warn, got %s", received.Level)
	}
}

func TestDispatch_HandlerError(t *testing.T) {
	n := New(0)
	n.Register(func(e Event) error {
		return errors.New("handler failed")
	})
	if err := n.Dispatch(443, "closed", LevelAlert); err == nil {
		t.Error("expected error from handler")
	}
}

func TestDispatch_Cooldown(t *testing.T) {
	n := New(5 * time.Second)
	calls := 0
	n.Register(func(e Event) error {
		calls++
		return nil
	})
	_ = n.Dispatch(80, "open", LevelInfo)
	_ = n.Dispatch(80, "open", LevelInfo) // should be suppressed
	if calls != 1 {
		t.Errorf("expected 1 call due to cooldown, got %d", calls)
	}
}

func TestDispatch_CooldownExpired(t *testing.T) {
	n := New(1 * time.Millisecond)
	calls := 0
	n.Register(func(e Event) error {
		calls++
		return nil
	})
	_ = n.Dispatch(80, "open", LevelInfo)
	time.Sleep(5 * time.Millisecond)
	_ = n.Dispatch(80, "open", LevelInfo)
	if calls != 2 {
		t.Errorf("expected 2 calls after cooldown expired, got %d", calls)
	}
}

func TestDispatch_MultipleHandlers(t *testing.T) {
	n := New(0)
	calls := 0
	for i := 0; i < 3; i++ {
		n.Register(func(e Event) error {
			calls++
			return nil
		})
	}
	_ = n.Dispatch(22, "open", LevelInfo)
	if calls != 3 {
		t.Errorf("expected 3 handler calls, got %d", calls)
	}
}
