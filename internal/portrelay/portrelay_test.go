package portrelay_test

import (
	"errors"
	"testing"

	"portwatch/internal/portrelay"
)

func TestRegister_Valid(t *testing.T) {
	r := portrelay.New()
	err := r.Register(portrelay.Target{
		Name:    "test",
		Handler: func(port int, event string) error { return nil },
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 target, got %d", r.Len())
	}
}

func TestRegister_EmptyName(t *testing.T) {
	r := portrelay.New()
	err := r.Register(portrelay.Target{Handler: func(int, string) error { return nil }})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilHandler(t *testing.T) {
	r := portrelay.New()
	err := r.Register(portrelay.Target{Name: "x"})
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestUnregister_RemovesTarget(t *testing.T) {
	r := portrelay.New()
	_ = r.Register(portrelay.Target{Name: "a", Handler: func(int, string) error { return nil }})
	r.Unregister("a")
	if r.Len() != 0 {
		t.Fatalf("expected 0 targets after unregister")
	}
}

func TestDispatch_CallsAllTargets(t *testing.T) {
	r := portrelay.New()
	called := map[string]bool{}
	for _, name := range []string{"a", "b"} {
		n := name
		_ = r.Register(portrelay.Target{
			Name: n,
			Handler: func(port int, event string) error {
				called[n] = true
				return nil
			},
		})
	}
	if err := r.Dispatch(8080, "open"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called["a"] || !called["b"] {
		t.Fatal("not all targets were called")
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	r := portrelay.New()
	_ = r.Register(portrelay.Target{
		Name:    "bad",
		Handler: func(int, string) error { return errors.New("fail") },
	})
	err := r.Dispatch(443, "close")
	if err == nil {
		t.Fatal("expected error from failing target")
	}
}

func TestDispatch_NoTargets(t *testing.T) {
	r := portrelay.New()
	if err := r.Dispatch(80, "open"); err != nil {
		t.Fatalf("unexpected error with no targets: %v", err)
	}
}
