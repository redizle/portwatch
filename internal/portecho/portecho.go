// Package portecho records the last observed echo (round-trip probe) result
// for each monitored port, tracking latency and whether the port responded.
package portecho

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portecho: invalid port number")

// Result holds the outcome of a single echo probe.
type Result struct {
	Port      int
	Latency   time.Duration
	Responded bool
	RecordedAt time.Time
}

// Echo stores the most recent probe result per port.
type Echo struct {
	mu      sync.RWMutex
	results map[int]Result
}

// New returns an initialised Echo store.
func New() *Echo {
	return &Echo{results: make(map[int]Result)}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Record saves the probe result for the given port.
func (e *Echo) Record(port int, latency time.Duration, responded bool) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.results[port] = Result{
		Port:       port,
		Latency:    latency,
		Responded:  responded,
		RecordedAt: time.Now(),
	}
	return nil
}

// Get returns the most recent probe result for the port.
// The second return value is false if no result has been recorded.
func (e *Echo) Get(port int) (Result, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	r, ok := e.results[port]
	return r, ok
}

// Clear removes the stored result for the given port.
func (e *Echo) Clear(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.results, port)
}

// Len returns the number of ports with recorded results.
func (e *Echo) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.results)
}
