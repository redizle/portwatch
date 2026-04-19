package porttag

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("porttag: invalid port number")

// Store holds a set of string tags per port.
type Store struct {
	mu   sync.RWMutex
	tags map[int]map[string]struct{}
}

// New returns an empty Store.
func New() *Store {
	return &Store{tags: make(map[int]map[string]struct{})}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Add attaches a tag to the given port.
func (s *Store) Add(port int, tag string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if tag == "" {
		return errors.New("porttag: tag must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tags[port] == nil {
		s.tags[port] = make(map[string]struct{})
	}
	s.tags[port][tag] = struct{}{}
	return nil
}

// Remove detaches a tag from the given port.
func (s *Store) Remove(port int, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tags[port], tag)
	if len(s.tags[port]) == 0 {
		delete(s.tags, port)
	}
}

// Get returns all tags for the given port.
func (s *Store) Get(port int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set := s.tags[port]
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	return out
}

// Has reports whether the given port has the given tag.
func (s *Store) Has(port int, tag string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.tags[port][tag]
	return ok
}

// Clear removes all tags for the given port.
func (s *Store) Clear(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tags, port)
}
