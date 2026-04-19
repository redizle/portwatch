// Package portbadge assigns short status badges to ports for display purposes.
package portbadge

import (
	"fmt"
	"sync"
)

// Badge represents a short display label for a port.
type Badge struct {
	Icon  string
	Label string
}

// Store holds badge assignments keyed by port number.
type Store struct {
	mu      sync.RWMutex
	badges  map[int]Badge
	default Badge
}

// New returns a new Store with a default badge.
func New() *Store {
	return &Store{
		badges:  make(map[int]Badge),
		default: Badge{Icon: "●", Label: "unknown"},
	}
}

func validPort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("portbadge: port %d out of range", p)
	}
	return nil
}

// Set assigns a badge to a port.
func (s *Store) Set(port int, icon, label string) error {
	if err := validPort(port); err != nil {
		return err
	}
	if label == "" {
		return fmt.Errorf("portbadge: label must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.badges[port] = Badge{Icon: icon, Label: label}
	return nil
}

// Get returns the badge for a port, falling back to the default.
func (s *Store) Get(port int) Badge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if b, ok := s.badges[port]; ok {
		return b
	}
	return s.default
}

// Remove deletes the badge for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.badges, port)
}

// All returns a copy of all assigned badges.
func (s *Store) All() map[int]Badge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]Badge, len(s.badges))
	for k, v := range s.badges {
		out[k] = v
	}
	return out
}
