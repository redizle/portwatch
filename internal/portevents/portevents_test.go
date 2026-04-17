package portevents_test

import (
	"sync"
	"testing"

	"portwatch/internal/portevents"
)

func TestPublish_NoHandlers(t *testing.T) {
	b := portevents.New()
	// should not panic
	b.Publish(portevents.Event{Port: 80, Type: portevents.EventOpened})
}

func TestSubscribe_And_Publish(t *testing.T) {
	b := portevents.New()
	var got portevents.Event
	b.Subscribe(portevents.EventOpened, func(e portevents.Event) {
		got = e
	})
	b.Publish(portevents.Event{Port: 443, Type: portevents.EventOpened})
	if got.Port != 443 {
		t.Fatalf("expected port 443, got %d", got.Port)
	}
}

func TestPublish_OnlyMatchingType(t *testing.T) {
	b := portevents.New()
	called := false
	b.Subscribe(portevents.EventClosed, func(e portevents.Event) {
		called = true
	})
	b.Publish(portevents.Event{Port: 80, Type: portevents.EventOpened})
	if called {
		t.Fatal("handler should not have been called for different event type")
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	b := portevents.New()
	b.Subscribe(portevents.EventOpened, func(e portevents.Event) {})
	b.Subscribe(portevents.EventOpened, func(e portevents.Event) {})
	if b.Len(portevents.EventOpened) != 2 {
		t.Fatalf("expected 2 handlers, got %d", b.Len(portevents.EventOpened))
	}
}

func TestPublish_MultipleHandlers(t *testing.T) {
	b := portevents.New()
	var mu sync.Mutex
	count := 0
	for i := 0; i < 3; i++ {
		b.Subscribe(portevents.EventChanged, func(e portevents.Event) {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}
	b.Publish(portevents.Event{Port: 22, Type: portevents.EventChanged})
	if count != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}
