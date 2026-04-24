// Package portprofile manages named profiles that group port configuration
// metadata such as label, owner, priority, and tags under a single identifier.
package portprofile

import (
	"errors"
	"fmt"
	"sync"
)

// Profile holds a named collection of metadata for a port.
type Profile struct {
	Name     string
	Label    string
	Owner    string
	Priority int
	Tags     []string
}

// Store maps port numbers to their assigned Profile.
type Store struct {
	mu       sync.RWMutex
	profiles map[int]Profile
}

// New returns an initialised Store.
func New() *Store {
	return &Store{profiles: make(map[int]Profile)}
}

func validPort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portprofile: invalid port %d", port)
	}
	return nil
}

// Set associates a Profile with the given port, overwriting any existing entry.
func (s *Store) Set(port int, p Profile) error {
	if err := validPort(port); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("portprofile: profile name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[port] = p
	return nil
}

// Get returns the Profile assigned to port and whether it was found.
func (s *Store) Get(port int) (Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profiles[port]
	return p, ok
}

// Remove deletes the profile for port. It is a no-op if no entry exists.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.profiles, port)
}

// All returns a copy of every port-to-profile mapping.
func (s *Store) All() map[int]Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]Profile, len(s.profiles))
	for k, v := range s.profiles {
		out[k] = v
	}
	return out
}
