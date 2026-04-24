// Package portrank assigns and tracks a numeric rank (priority score) to ports
// based on activity frequency, alert count, and manual overrides.
package portrank

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portrank: invalid port number")

// Entry holds the computed rank and contributing factors for a port.
type Entry struct {
	Port      int
	Score     int
	Override  bool
}

// Ranker stores rank scores for ports.
type Ranker struct {
	mu      sync.RWMutex
	scores  map[int]int
	override map[int]int
}

// New returns an initialised Ranker.
func New() *Ranker {
	return &Ranker{
		scores:   make(map[int]int),
		override: make(map[int]int),
	}
}

func validPort(p int) bool {
	return p >= 1 && p <= 65535
}

// Add increments the activity score for a port by delta.
func (r *Ranker) Add(port, delta int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scores[port] += delta
	return nil
}

// SetOverride pins a port to a fixed score, ignoring activity accumulation.
func (r *Ranker) SetOverride(port, score int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.override[port] = score
	return nil
}

// ClearOverride removes a manual override, restoring activity-based scoring.
func (r *Ranker) ClearOverride(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.override, port)
}

// Get returns the effective Entry for a port.
func (r *Ranker) Get(port int) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.override[port]; ok {
		return Entry{Port: port, Score: v, Override: true}, true
	}
	if v, ok := r.scores[port]; ok {
		return Entry{Port: port, Score: v}, true
	}
	return Entry{}, false
}

// Reset clears both the activity score and override for a port.
func (r *Ranker) Reset(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.scores, port)
	delete(r.override, port)
}

// Len returns the number of ports with any recorded score or override.
func (r *Ranker) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	seen := make(map[int]struct{}, len(r.scores)+len(r.override))
	for p := range r.scores {
		seen[p] = struct{}{}
	}
	for p := range r.override {
		seen[p] = struct{}{}
	}
	return len(seen)
}
