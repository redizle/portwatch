// Package portquota tracks connection attempt counts against configurable
// thresholds, allowing callers to detect when a port exceeds its quota.
package portquota

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portquota: port must be between 1 and 65535")

// ErrInvalidQuota is returned when a quota limit is zero or negative.
var ErrInvalidQuota = errors.New("portquota: quota must be greater than zero")

// Entry holds the current hit count and maximum quota for a port.
type Entry struct {
	Hits  int
	Limit int
}

// Exceeded reports whether the hit count has reached or surpassed the limit.
func (e Entry) Exceeded() bool { return e.Hits >= e.Limit }

// Quota manages per-port hit counts and limits.
type Quota struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an initialised Quota.
func New() *Quota {
	return &Quota{entries: make(map[int]*Entry)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set configures a quota limit for the given port.
func (q *Quota) Set(port, limit int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if limit <= 0 {
		return ErrInvalidQuota
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if e, ok := q.entries[port]; ok {
		e.Limit = limit
	} else {
		q.entries[port] = &Entry{Limit: limit}
	}
	return nil
}

// Inc increments the hit counter for a port. Returns ErrInvalidPort for bad ports.
func (q *Quota) Inc(port int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if e, ok := q.entries[port]; ok {
		e.Hits++
	}
	return nil
}

// Get returns the Entry for a port and whether it exists.
func (q *Quota) Get(port int) (Entry, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	e, ok := q.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Reset clears the hit count for a port without removing its limit.
func (q *Quota) Reset(port int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if e, ok := q.entries[port]; ok {
		e.Hits = 0
	}
}

// Remove deletes all quota data for a port.
func (q *Quota) Remove(port int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.entries, port)
}
