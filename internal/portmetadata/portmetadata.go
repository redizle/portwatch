// Package portmetadata provides a generic key-value metadata store for ports.
package portmetadata

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("invalid port: must be between 1 and 65535")

// ErrEmptyKey is returned when a metadata key is empty.
var ErrEmptyKey = errors.New("metadata key must not be empty")

// Store holds arbitrary string metadata keyed by port and field name.
type Store struct {
	mu   sync.RWMutex
	data map[int]map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[int]map[string]string)}
}

func validPort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%w: %d", ErrInvalidPort, port)
	}
	return nil
}

// Set stores a metadata value for the given port and key.
func (s *Store) Set(port int, key, value string) error {
	if err := validPort(port); err != nil {
		return err
	}
	if key == "" {
		return ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[port] == nil {
		s.data[port] = make(map[string]string)
	}
	s.data[port][key] = value
	return nil
}

// Get retrieves a metadata value. Returns ("", false) if not found.
func (s *Store) Get(port int, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.data[port]
	if !ok {
		return "", false
	}
	v, ok := m[key]
	return v, ok
}

// Delete removes a single key for a port.
func (s *Store) Delete(port int, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.data[port]; ok {
		delete(m, key)
		if len(m) == 0 {
			delete(s.data, port)
		}
	}
}

// All returns a copy of all metadata for the given port.
func (s *Store) All(port int) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string)
	for k, v := range s.data[port] {
		out[k] = v
	}
	return out
}

// Clear removes all metadata for the given port.
func (s *Store) Clear(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, port)
}
