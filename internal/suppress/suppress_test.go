package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func TestSuppress_IsSuppressed(t *testing.T) {
	s := suppress.New()
	s.Suppress(8080, "maintenance", 1*time.Hour)
	if !s.IsSuppressed(8080) {
		t.Fatal("expected port 8080 to be suppressed")
	}
}

func TestSuppress_NotSuppressed_UnknownPort(t *testing.T) {
	s := suppress.New()
	if s.IsSuppressed(9999) {
		t.Fatal("expected port 9999 to not be suppressed")
	}
}

func TestSuppress_Expired(t *testing.T) {
	s := suppress.New()
	s.Suppress(3000, "test", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if s.IsSuppressed(3000) {
		t.Fatal("expected suppression to have expired")
	}
}

func TestLift_RemovesSuppression(t *testing.T) {
	s := suppress.New()
	s.Suppress(443, "deploy", 1*time.Hour)
	s.Lift(443)
	if s.IsSuppressed(443) {
		t.Fatal("expected suppression to be lifted")
	}
}

// TestLift_NoopOnUnknownPort ensures Lift does not panic on an unknown port.
func TestLift_NoopOnUnknownPort(t *testing.T) {
	s := suppress.New()
	s.Lift(1234) // should not panic
}

func TestActive_ReturnsOnlyLive(t *testing.T) {
	s := suppress.New()
	s.Suppress(80, "a", 1*time.Hour)
	s.Suppress(81, "b", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	actives := s.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 active suppression, got %d", len(actives))
	}
	if actives[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", actives[0].Port)
	}
}

func TestActive_Empty(t *testing.T) {
	s := suppress.New()
	if len(s.Active()) != 0 {
		t.Fatal("expected no active suppressions")
	}
}
