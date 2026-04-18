// Package portcomment allows attaching freeform comments to ports.
package portcomment

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portcomment: port must be between 1 and 65535")

// Store holds comments keyed by port number.
type Store struct {
	mu       sync.RWMutex
	comments map[int]string
}

// New returns an empty comment Store.
func New() *Store {
	return &Store{comments: make(map[int]string)}
}

func validate(port int) error {
	if port < 1 || port > 65535 {
		return ErrInvalidPort
	}
	return nil
}

// Set attaches a comment to the given port.
func (s *Store) Set(port int, comment string) error {
	if err := validate(port); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[port] = comment
	return nil
}

// Get returns the comment for a port and whether it exists.
func (s *Store) Get(port int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.comments[port]
	return c, ok
}

// Remove deletes the comment for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.comments, port)
}

// All returns a snapshot of all port comments.
func (s *Store) All() map[int]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]string, len(s.comments))
	for k, v := range s.comments {
		out[k] = v
	}
	return out
}
