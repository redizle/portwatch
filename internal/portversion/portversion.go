// Package portversion tracks a version/revision string associated with
// a monitored port — useful for noting which service version was observed.
package portversion

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portversion: port must be between 1 and 65535")

// ErrEmptyVersion is returned when an empty version string is set.
var ErrEmptyVersion = errors.New("portversion: version must not be empty")

// Entry holds the version string for a port.
type Entry struct {
	Port    int
	Version string
}

// Store maps ports to version strings.
type Store struct {
	mu      sync.RWMutex
	entries map[int]string
}

// New creates an empty Store.
func New() *Store {
	return &Store{entries: make(map[int]string)}
}

func validPort(p int) bool {
	return p >= 1 && p <= 65535
}

// Set assigns a version string to a port.
func (s *Store) Set(port int, version string) error {
	if !validPort(port) {
		return fmt.Errorf("%w: got %d", ErrInvalidPort, port)
	}
	if version == "" {
		return ErrEmptyVersion
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = version
	return nil
}

// Get returns the version for a port and whether it was found.
func (s *Store) Get(port int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.entries[port]
	return v, ok
}

// Remove deletes the version entry for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// All returns a snapshot of all entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for p, v := range s.entries {
		out = append(out, Entry{Port: p, Version: v})
	}
	return out
}
