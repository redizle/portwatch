// Package portseq tracks the order in which ports were first observed open.
package portseq

import (
	"errors"
	"sync"
)

// Entry holds sequence metadata for a port.
type Entry struct {
	Port     int
	Sequence int
}

// Sequencer assigns monotonically increasing sequence numbers to ports.
type Sequencer struct {
	mu      sync.Mutex
	counter int
	ports   map[int]Entry
}

// New returns an empty Sequencer.
func New() *Sequencer {
	return &Sequencer{ports: make(map[int]Entry)}
}

// Record assigns the next sequence number to port if not already recorded.
// Returns the entry and whether it was newly assigned.
func (s *Sequencer) Record(port int) (Entry, bool, error) {
	if port < 1 || port > 65535 {
		return Entry{}, false, errors.New("portseq: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.ports[port]; ok {
		return e, false, nil
	}
	s.counter++
	e := Entry{Port: port, Sequence: s.counter}
	s.ports[port] = e
	return e, true, nil
}

// Get returns the entry for a port, or false if not recorded.
func (s *Sequencer) Get(port int) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.ports[port]
	return e, ok
}

// Reset clears all recorded ports and resets the counter.
func (s *Sequencer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = 0
	s.ports = make(map[int]Entry)
}

// Len returns the number of recorded ports.
func (s *Sequencer) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.ports)
}
