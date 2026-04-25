// Package portweight assigns and tracks numeric weights to monitored ports.
// Weights influence prioritization in reporting and alerting pipelines.
package portweight

import (
	"errors"
	"sync"
)

const defaultWeight = 1

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portweight: port must be between 1 and 65535")

// ErrInvalidWeight is returned when a weight is not positive.
var ErrInvalidWeight = errors.New("portweight: weight must be greater than zero")

// Entry holds the weight assigned to a port.
type Entry struct {
	Port   int
	Weight int
}

// Weights stores per-port weight values.
type Weights struct {
	mu      sync.RWMutex
	weights map[int]int
}

// New returns an initialised Weights store.
func New() *Weights {
	return &Weights{weights: make(map[int]int)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set assigns weight w to port p.
func (wt *Weights) Set(port, w int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if w <= 0 {
		return ErrInvalidWeight
	}
	wt.mu.Lock()
	defer wt.mu.Unlock()
	wt.weights[port] = w
	return nil
}

// Get returns the weight for port p, falling back to defaultWeight.
func (wt *Weights) Get(port int) int {
	wt.mu.RLock()
	defer wt.mu.RUnlock()
	if v, ok := wt.weights[port]; ok {
		return v
	}
	return defaultWeight
}

// Remove deletes any explicit weight for port p.
func (wt *Weights) Remove(port int) {
	wt.mu.Lock()
	defer wt.mu.Unlock()
	delete(wt.weights, port)
}

// All returns a snapshot of all explicitly set entries.
func (wt *Weights) All() []Entry {
	wt.mu.RLock()
	defer wt.mu.RUnlock()
	out := make([]Entry, 0, len(wt.weights))
	for p, w := range wt.weights {
		out = append(out, Entry{Port: p, Weight: w})
	}
	return out
}

// Len returns the number of explicitly weighted ports.
func (wt *Weights) Len() int {
	wt.mu.RLock()
	defer wt.mu.RUnlock()
	return len(wt.weights)
}
