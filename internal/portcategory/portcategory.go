// Package portcategory assigns user-defined categories to ports.
package portcategory

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("port must be between 1 and 65535")

// ErrEmptyCategory is returned when an empty category string is provided.
var ErrEmptyCategory = errors.New("category must not be empty")

// Store holds category assignments for ports.
type Store struct {
	mu         sync.RWMutex
	categories map[int]string
}

// New creates a new Store.
func New() *Store {
	return &Store{categories: make(map[int]string)}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Set assigns a category to a port.
func (s *Store) Set(port int, category string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if category == "" {
		return ErrEmptyCategory
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.categories[port] = category
	return nil
}

// Get returns the category for a port and whether it was found.
func (s *Store) Get(port int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.categories[port]
	return c, ok
}

// Remove deletes the category assignment for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.categories, port)
}

// All returns a copy of all port-to-category mappings.
func (s *Store) All() map[int]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]string, len(s.categories))
	for k, v := range s.categories {
		out[k] = v
	}
	return out
}

// ByCategory returns all ports assigned to the given category.
func (s *Store) ByCategory(category string) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ports []int
	for port, cat := range s.categories {
		if cat == category {
			ports = append(ports, port)
		}
	}
	return ports
}
