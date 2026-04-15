package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstCallPasses(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	th.Allow()
	if th.Allow() {
		t.Fatal("expected second immediate call to be blocked")
	}
	if th.Skipped() != 1 {
		t.Fatalf("expected 1 skipped, got %d", th.Skipped())
	}
}

func TestAllow_PassesAfterInterval(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow()
	time.Sleep(30 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected call to pass after interval elapsed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected call to pass after reset")
	}
}

func TestSetInterval_UpdatesThreshold(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	th.Allow()
	th.SetInterval(1 * time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected call to pass after interval was shortened")
	}
}

func TestSkipped_CountsCorrectly(t *testing.T) {
	th := throttle.New(1 * time.Second)
	th.Allow()
	for i := 0; i < 5; i++ {
		th.Allow()
	}
	if th.Skipped() != 5 {
		t.Fatalf("expected 5 skipped, got %d", th.Skipped())
	}
}
