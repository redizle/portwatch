package portcensor

import (
	"testing"
)

func TestRedact_And_IsCensored(t *testing.T) {
	c := New()
	if err := c.Redact(443, "TLS traffic"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.IsCensored(443) {
		t.Error("expected port 443 to be censored")
	}
}

func TestRedact_InvalidPort(t *testing.T) {
	c := New()
	if err := c.Redact(0, "bad"); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
	if err := c.Redact(70000, "bad"); err != ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestRedact_EmptyReason(t *testing.T) {
	c := New()
	if err := c.Redact(80, ""); err != ErrEmptyReason {
		t.Errorf("expected ErrEmptyReason, got %v", err)
	}
}

func TestIsCensored_Missing(t *testing.T) {
	c := New()
	if c.IsCensored(8080) {
		t.Error("expected port 8080 to not be censored")
	}
}

func TestLift_RemovesCensor(t *testing.T) {
	c := New()
	_ = c.Redact(22, "ssh")
	c.Lift(22)
	if c.IsCensored(22) {
		t.Error("expected port 22 to be lifted")
	}
}

func TestLift_NoopOnMissing(t *testing.T) {
	c := New()
	c.Lift(9999) // should not panic
}

func TestGet_Found(t *testing.T) {
	c := New()
	_ = c.Redact(5432, "postgres")
	e, ok := c.Get(5432)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Port != 5432 || e.Reason != "postgres" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestGet_Missing(t *testing.T) {
	c := New()
	_, ok := c.Get(1234)
	if ok {
		t.Error("expected missing port to return false")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	c := New()
	_ = c.Redact(80, "http")
	_ = c.Redact(443, "https")
	all := c.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestLen(t *testing.T) {
	c := New()
	if c.Len() != 0 {
		t.Error("expected empty censor")
	}
	_ = c.Redact(3306, "mysql")
	if c.Len() != 1 {
		t.Errorf("expected 1, got %d", c.Len())
	}
}
