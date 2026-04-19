// Package portwindow tracks time-based observation windows for ports,
// recording when a port was first and last seen within a rolling window.
package portwindow

import (
	"fmt"
	"sync"
	"time"
)

// Window holds the observation window for a single port.
type Window struct {
	Port     int
	First    time.Time
	Last     time.Time
	Hits     int
	Duration time.Duration
}

// Tracker manages rolling observation windows per port.
type Tracker struct {
	mu      sync.Mutex
	windows map[int]*Window
	size    time.Duration
}

// New creates a Tracker with the given window size.
func New(size time.Duration) (*Tracker, error) {
	if size <= 0 {
		return nil, fmt.Errorf("portwindow: window size must be positive")
	}
	return &Tracker{
		windows: make(map[int]*Window),
		size:    size,
	}, nil
}

// Observe records an observation for the given port at now.
func (t *Tracker) Observe(port int, now time.Time) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portwindow: invalid port %d", port)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	w, ok := t.windows[port]
	if !ok || now.Sub(w.First) > t.size {
		t.windows[port] = &Window{Port: port, First: now, Last: now, Hits: 1}
		return nil
	}
	w.Last = now
	w.Hits++
	w.Duration = w.Last.Sub(w.First)
	return nil
}

// Get returns the current window for a port, or false if not present.
func (t *Tracker) Get(port int) (Window, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	w, ok := t.windows[port]
	if !ok {
		return Window{}, false
	}
	return *w, true
}

// Reset clears the window for a port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.windows, port)
}

// Len returns the number of tracked ports.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.windows)
}
