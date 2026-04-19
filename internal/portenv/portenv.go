// Package portenv associates an environment tag (e.g. "prod", "staging", "dev")
// with a port number for richer reporting and alerting context.
package portenv

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portenv: port must be between 1 and 65535")

// ErrEmptyEnv is returned when an empty environment string is provided.
var ErrEmptyEnv = errors.New("portenv: environment must not be empty")

// Entry holds the environment tag for a port.
type Entry struct {
	Port int
	Env  string
}

// Store maps ports to environment tags.
type Store struct {
	mu      sync.RWMutex
	entries map[int]string
	default_ string
}

// New creates a new Store with the given default environment fallback.
func New(defaultEnv string) *Store {
	return &Store{
		entries:  make(map[int]string),
		default_: defaultEnv,
	}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set assigns an environment tag to a port.
func (s *Store) Set(port int, env string) error {
	if !validPort(port) {
		return fmt.Errorf("%w: got %d", ErrInvalidPort, port)
	}
	if env == "" {
		return ErrEmptyEnv
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = env
	return nil
}

// Get returns the environment for the given port, falling back to the default.
func (s *Store) Get(port int) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if env, ok := s.entries[port]; ok {
		return env
	}
	return s.default_
}

// Remove deletes the environment override for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// All returns a snapshot of all explicitly set entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for p, e := range s.entries {
		out = append(out, Entry{Port: p, Env: e})
	}
	return out
}
