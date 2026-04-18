// Package portannot provides structured annotation storage for monitored ports.
package portannot

import (
	"errors"
	"sync"
)

// Annotation holds a key-value pair attached to a port.
type Annotation struct {
	Key   string
	Value string
}

// Store manages annotations per port.
type Store struct {
	mu   sync.RWMutex
	data map[int]map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[int]map[string]string)}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Set adds or updates an annotation key for the given port.
func (s *Store) Set(port int, key, value string) error {
	if !validPort(port) {
		return errors.New("portannot: port out of range")
	}
	if key == "" {
		return errors.New("portannot: key must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[port] == nil {
		s.data[port] = make(map[string]string)
	}
	s.data[port][key] = value
	return nil
}

// Get returns the value for a key on a port and whether it was found.
func (s *Store) Get(port int, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[port][key]
	return v, ok
}

// Remove deletes a single annotation key from a port.
func (s *Store) Remove(port int, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data[port], key)
}

// All returns a copy of all annotations for a port.
func (s *Store) All(port int) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.data[port]))
	for k, v := range s.data[port] {
		out[k] = v
	}
	return out
}

// Clear removes all annotations for a port.
func (s *Store) Clear(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, port)
}
