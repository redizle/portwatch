package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallPasses(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	if !l.Allow("port:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("port:8080")
	if l.Allow("port:8080") {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestAllow_PassesAfterInterval(t *testing.T) {
	l := ratelimit.New(30 * time.Millisecond)
	l.Allow("port:9090")
	time.Sleep(40 * time.Millisecond)
	if !l.Allow("port:9090") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("port:8080")
	if !l.Allow("port:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("port:8080")
	l.Reset("port:8080")
	if !l.Allow("port:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("port:8080")
	l.Allow("port:9090")
	l.Flush()
	if !l.Allow("port:8080") {
		t.Fatal("expected allow after flush for port:8080")
	}
	if !l.Allow("port:9090") {
		t.Fatal("expected allow after flush for port:9090")
	}
}
