// Package watchlist manages a user-defined set of ports to monitor with priority.
package watchlist

import (
	"fmt"
	"sync"
)

// Priority represents how important a port is to the user.
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 2
	PriorityHigh   Priority = 3
)

// Entry holds metadata about a watched port.
type Entry struct {
	Port     int
	Label    string
	Priority Priority
}

// Watchlist stores a prioritised set of ports.
type Watchlist struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{entries: make(map[int]Entry)}
}

// Add registers a port with the given label and priority.
// Returns an error if the port number is out of range.
func (w *Watchlist) Add(port int, label string, priority Priority) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("watchlist: port %d out of range", port)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[port] = Entry{Port: port, Label: label, Priority: priority}
	return nil
}

// Remove deletes a port from the watchlist. No-op if not present.
func (w *Watchlist) Remove(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, port)
}

// Get returns the Entry for a port and whether it was found.
func (w *Watchlist) Get(port int) (Entry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	e, ok := w.entries[port]
	return e, ok
}

// All returns a snapshot of all entries.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// Contains reports whether a port is in the watchlist.
func (w *Watchlist) Contains(port int) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.entries[port]
	return ok
}
