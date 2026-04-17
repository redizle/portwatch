// Package portevents provides a simple event bus for port state changes.
package portevents

import "sync"

// EventType describes the kind of port event.
type EventType string

const (
	EventOpened  EventType = "opened"
	EventClosed  EventType = "closed"
	EventChanged EventType = "changed"
)

// Event represents a single port state change.
type Event struct {
	Port  int
	Type  EventType
	Extra map[string]string
}

// Handler is a function that handles a port event.
type Handler func(e Event)

// Bus dispatches port events to registered handlers.
type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

// New returns a new Bus.
func New() *Bus {
	return &Bus{handlers: make(map[EventType][]Handler)}
}

// Subscribe registers a handler for the given event type.
func (b *Bus) Subscribe(t EventType, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

// Publish dispatches an event to all matching handlers.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers[e.Type] {
		h(e)
	}
}

// Len returns the number of handlers registered for a given type.
func (b *Bus) Len(t EventType) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[t])
}
