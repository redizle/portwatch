// Package portalias maps ports to human-friendly alias names.
package portalias

import (
	"errors"
	"sync"
)

var (
	ErrInvalidPort = errors.New("port must be between 1 and 65535")
	ErrEmptyAlias  = errors.New("alias must not be empty")
)

// Alias holds an alias entry for a port.
type Alias struct {
	Port  int
	Label string
}

// Store manages port-to-alias mappings.
type Store struct {
	mu      sync.RWMutex
	aliases map[int]string
}

// New returns an empty alias Store.
func New() *Store {
	return &Store{aliases: make(map[int]string)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set assigns an alias to a port.
func (s *Store) Set(port int, label string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if label == "" {
		return ErrEmptyAlias
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.aliases[port] = label
	return nil
}

// Get returns the alias for a port and whether it was found.
func (s *Store) Get(port int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.aliases[port]
	return v, ok
}

// Remove deletes the alias for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.aliases, port)
}

// All returns a snapshot of all alias entries.
func (s *Store) All() []Alias {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Alias, 0, len(s.aliases))
	for p, l := range s.aliases {
		out = append(out, Alias{Port: p, Label: l})
	}
	return out
}
