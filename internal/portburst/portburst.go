// Package portburst tracks short-term burst activity for ports,
// flagging ports that exceed a hit threshold within a rolling time window.
package portburst

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portburst: port must be between 1 and 65535")

// ErrInvalidThreshold is returned when the burst threshold is zero or negative.
var ErrInvalidThreshold = errors.New("portburst: threshold must be greater than zero")

// entry holds timestamped hit records for a single port.
type entry struct {
	hits   []time.Time
}

// Tracker records hits per port and reports bursts.
type Tracker struct {
	mu        sync.Mutex
	entries   map[int]*entry
	window    time.Duration
	threshold int
}

// New creates a Tracker that flags a port as bursting when it receives
// at least threshold hits within the given rolling window.
func New(window time.Duration, threshold int) (*Tracker, error) {
	if threshold <= 0 {
		return nil, ErrInvalidThreshold
	}
	return &Tracker{
		entries:   make(map[int]*entry),
		window:    window,
		threshold: threshold,
	}, nil
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Record registers a hit for the given port at the current time.
func (t *Tracker) Record(port int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		e = &entry{}
		t.entries[port] = e
	}
	e.hits = append(e.hits, time.Now())
	return nil
}

// IsBursting returns true if the port has received at least threshold hits
// within the configured window. Returns false for unknown ports.
func (t *Tracker) IsBursting(port int) bool {
	if !validPort(port) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		return false
	}
	cutoff := time.Now().Add(-t.window)
	count := 0
	for _, ts := range e.hits {
		if ts.After(cutoff) {
			count++
		}
	}
	return count >= t.threshold
}

// Reset clears all recorded hits for the given port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// HitCount returns the number of hits within the current window for a port.
func (t *Tracker) HitCount(port int) int {
	if !validPort(port) {
		return 0
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		return 0
	}
	cutoff := time.Now().Add(-t.window)
	count := 0
	for _, ts := range e.hits {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
