// Package snapshot captures and compares port state snapshots
// to detect changes between scan cycles.
package snapshot

import (
	"sync"
	"time"
)

// PortState represents the observed state of a single port.
type PortState struct {
	Port      int
	Open      bool
	FirstSeen time.Time
	LastSeen  time.Time
}

// Snapshot holds a point-in-time view of all monitored ports.
type Snapshot struct {
	mu    sync.RWMutex
	ports map[int]PortState
	Taken time.Time
}

// New returns an empty Snapshot.
func New() *Snapshot {
	return &Snapshot{
		ports: make(map[int]PortState),
		Taken: time.Now(),
	}
}

// Set records the state of a port in the snapshot.
func (s *Snapshot) Set(port int, open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.ports[port]
	now := time.Now()

	if !ok {
		s.ports[port] = PortState{
			Port:      port,
			Open:      open,
			FirstSeen: now,
			LastSeen:  now,
		}
		return
	}

	existing.Open = open
	existing.LastSeen = now
	s.ports[port] = existing
}

// Get returns the PortState for a given port and whether it exists.
func (s *Snapshot) Get(port int) (PortState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ps, ok := s.ports[port]
	return ps, ok
}

// All returns a copy of all recorded port states.
func (s *Snapshot) All() []PortState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]PortState, 0, len(s.ports))
	for _, ps := range s.ports {
		result = append(result, ps)
	}
	return result
}

// Diff compares this snapshot against a previous one and returns
// ports whose open/closed status changed.
func (s *Snapshot) Diff(prev *Snapshot) []PortState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var changed []PortState
	for port, cur := range s.ports {
		if prev == nil {
			if cur.Open {
				changed = append(changed, cur)
			}
			continue
		}
		prev.mu.RLock()
		prevState, ok := prev.ports[port]
		prev.mu.RUnlock()
		if !ok || prevState.Open != cur.Open {
			changed = append(changed, cur)
		}
	}
	return changed
}
