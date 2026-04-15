package state

import (
	"sync"
	"time"
)

// PortStatus represents the current known status of a port.
type PortStatus struct {
	Port     int
	Open     bool
	LastSeen time.Time
	FirstSeen time.Time
}

// Store holds the last known state of all monitored ports.
type Store struct {
	mu      sync.RWMutex
	ports   map[int]*PortStatus
}

// New creates a new state Store.
func New() *Store {
	return &Store{
		ports: make(map[int]*PortStatus),
	}
}

// Update sets the current status for a port and returns whether the state changed.
func (s *Store) Update(port int, open bool) (changed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	prev, exists := s.ports[port]

	if !exists {
		s.ports[port] = &PortStatus{
			Port:      port,
			Open:      open,
			LastSeen:  now,
			FirstSeen: now,
		}
		return true
	}

	if prev.Open != open {
		prev.Open = open
		prev.LastSeen = now
		return true
	}

	prev.LastSeen = now
	return false
}

// Get returns the status for a given port, and whether it exists.
func (s *Store) Get(port int) (*PortStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ps, ok := s.ports[port]
	if !ok {
		return nil, false
	}
	copy := *ps
	return &copy, true
}

// Snapshot returns a copy of all port statuses.
func (s *Store) Snapshot() []PortStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]PortStatus, 0, len(s.ports))
	for _, ps := range s.ports {
		result = append(result, *ps)
	}
	return result
}
