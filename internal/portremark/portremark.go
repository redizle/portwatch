// Package portremark allows attaching short timestamped remarks to ports.
package portremark

import (
	"errors"
	"sync"
	"time"
)

// Remark holds a note and when it was created.
type Remark struct {
	Text      string
	CreatedAt time.Time
}

// Store manages remarks keyed by port number.
type Store struct {
	mu      sync.RWMutex
	remarks map[int]Remark
}

// New returns an empty Store.
func New() *Store {
	return &Store{remarks: make(map[int]Remark)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set attaches a remark to the given port.
func (s *Store) Set(port int, text string) error {
	if !validPort(port) {
		return errors.New("portremark: port out of range")
	}
	if text == "" {
		return errors.New("portremark: remark text must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.remarks[port] = Remark{Text: text, CreatedAt: time.Now()}
	return nil
}

// Get returns the remark for a port and whether it exists.
func (s *Store) Get(port int) (Remark, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.remarks[port]
	return r, ok
}

// Remove deletes the remark for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.remarks, port)
}

// All returns a copy of all remarks.
func (s *Store) All() map[int]Remark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]Remark, len(s.remarks))
	for k, v := range s.remarks {
		out[k] = v
	}
	return out
}
