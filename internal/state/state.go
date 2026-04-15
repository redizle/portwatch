package state

import (
	"sync"
	"time"
)

// PortEntry holds the last known status of a port.
type PortEntry struct {
	Open      bool
	LastSeen  time.Time
	ChangedAt time.Time
}

// State tracks the current open/closed status of monitored ports.
type State struct {
	mu    sync.RWMutex
	ports map[int]PortEntry
}

// New creates an empty State.
func New() *State {
	return &State{ports: make(map[int]PortEntry)}
}

// Update records a new status for the given port.
// Returns true if the status changed from the previous value.
func (s *State) Update(port int, open bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	prev, exists := s.ports[port]
	changed := !exists || prev.Open != open
	entry := PortEntry{
		Open:      open,
		LastSeen:  now,
		ChangedAt: prev.ChangedAt,
	}
	if changed {
		entry.ChangedAt = now
	}
	s.ports[port] = entry
	return changed
}

// Get returns the entry for a port and whether it exists.
func (s *State) Get(port int) (PortEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.ports[port]
	return e, ok
}

// All returns a snapshot copy of all tracked port entries.
func (s *State) All() map[int]PortEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]PortEntry, len(s.ports))
	for k, v := range s.ports {
		out[k] = v
	}
	return out
}

// Delete removes a port from tracking.
func (s *State) Delete(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ports, port)
}
