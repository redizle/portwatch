// Package notify provides threshold-based notification logic for port state changes.
package notify

import (
	"fmt"
	"sync"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelAlert Level = "alert"
)

// Event represents a notification event for a port change.
type Event struct {
	Port      int
	Status    string
	Level     Level
	Message   string
	Timestamp time.Time
}

// Handler is a function that receives a notification event.
type Handler func(Event) error

// Notifier dispatches events to registered handlers.
type Notifier struct {
	mu       sync.RWMutex
	handlers []Handler
	cooldown time.Duration
	lastSent map[int]time.Time
}

// New creates a Notifier with the given cooldown between repeated alerts per port.
func New(cooldown time.Duration) *Notifier {
	return &Notifier{
		cooldown: cooldown,
		lastSent: make(map[int]time.Time),
	}
}

// Register adds a handler to the notifier.
func (n *Notifier) Register(h Handler) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.handlers = append(n.handlers, h)
}

// Dispatch sends an event to all registered handlers, respecting cooldown.
func (n *Notifier) Dispatch(port int, status string, level Level) error {
	n.mu.Lock()
	if last, ok := n.lastSent[port]; ok && time.Since(last) < n.cooldown {
		n.mu.Unlock()
		return nil
	}
	n.lastSent[port] = time.Now()
	handlers := make([]Handler, len(n.handlers))
	copy(handlers, n.handlers)
	n.mu.Unlock()

	event := Event{
		Port:      port,
		Status:    status,
		Level:     level,
		Message:   fmt.Sprintf("port %d is %s", port, status),
		Timestamp: time.Now(),
	}

	var firstErr error
	for _, h := range handlers {
		if err := h(event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
