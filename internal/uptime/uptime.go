// Package uptime tracks how long individual ports have been continuously open.
package uptime

import (
	"sync"
	"time"
)

// Tracker records when each port was first seen open and calculates uptime.
type Tracker struct {
	mu    sync.Mutex
	opened map[int]time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		opened: make(map[int]time.Time),
	}
}

// MarkOpen records the time a port became open. If already tracked, it is a no-op.
func (t *Tracker) MarkOpen(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.opened[port]; !ok {
		t.opened[port] = time.Now()
	}
}

// MarkClosed removes a port from tracking.
func (t *Tracker) MarkClosed(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.opened, port)
}

// Uptime returns how long a port has been open. Returns 0 and false if not tracked.
func (t *Tracker) Uptime(port int) (time.Duration, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	open, ok := t.opened[port]
	if !ok {
		return 0, false
	}
	return time.Since(open), true
}

// OpenedAt returns the time a port was first seen open.
func (t *Tracker) OpenedAt(port int) (time.Time, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t2, ok := t.opened[port]
	return t2, ok
}

// All returns a snapshot of all tracked ports and their open times.
func (t *Tracker) All() map[int]time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[int]time.Time, len(t.opened))
	for k, v := range t.opened {
		out[k] = v
	}
	return out
}
