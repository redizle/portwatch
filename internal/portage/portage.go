// Package portage tracks how long a port has been in its current state.
package portage

import (
	"errors"
	"sync"
	"time"
)

// Entry holds age information for a single port.
type Entry struct {
	Port      int
	Since     time.Time
	LastSeen  time.Time
}

// Age returns how long the port has been in its current state.
func (e Entry) Age() time.Duration {
	return time.Since(e.Since)
}

// Tracker records when each port entered its current state.
type Tracker struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[int]Entry)}
}

var errInvalidPort = errors.New("portage: port must be between 1 and 65535")

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Mark records or refreshes the state-entry timestamp for a port.
// If the port is already tracked its LastSeen is updated; Since is only
// set on first observation.
func (t *Tracker) Mark(port int) error {
	if !validPort(port) {
		return errInvalidPort
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if e, ok := t.entries[port]; ok {
		e.LastSeen = now
		t.entries[port] = e
	} else {
		t.entries[port] = Entry{Port: port, Since: now, LastSeen: now}
	}
	return nil
}

// Reset clears the tracked entry so the next Mark starts a fresh Since.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Get returns the Entry for a port and whether it was found.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	return e, ok
}

// All returns a snapshot of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
