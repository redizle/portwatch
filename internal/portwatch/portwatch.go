// Package portwatch combines enrichment, labeling, and state tracking
// into a unified per-port watch record.
package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds a consolidated view of a watched port.
type Entry struct {
	Port      int
	Label     string
	Owner     string
	Open      bool
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
}

// Watcher tracks port entries by port number.
type Watcher struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an initialised Watcher.
func New() *Watcher {
	return &Watcher{entries: make(map[int]*Entry)}
}

// Touch records a port observation, creating the entry if needed.
func (w *Watcher) Touch(port int, open bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	e, ok := w.entries[port]
	if !ok {
		e = &Entry{Port: port, FirstSeen: now}
		w.entries[port] = e
	}
	e.Open = open
	e.LastSeen = now
	e.SeenCount++
}

// SetLabel attaches a human-readable label to a port entry.
func (w *Watcher) SetLabel(port int, label string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portwatch: invalid port %d", port)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.entries[port]; !ok {
		w.entries[port] = &Entry{Port: port}
	}
	w.entries[port].Label = label
	return nil
}

// SetOwner attaches an owner string to a port entry.
func (w *Watcher) SetOwner(port int, owner string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portwatch: invalid port %d", port)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.entries[port]; !ok {
		w.entries[port] = &Entry{Port: port}
	}
	w.entries[port].Owner = owner
	return nil
}

// Get returns a copy of the entry for the given port.
func (w *Watcher) Get(port int) (Entry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	e, ok := w.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a slice of all current entries.
func (w *Watcher) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, *e)
	}
	return out
}
