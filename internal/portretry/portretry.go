// Package portretry tracks retry attempts for ports that have failed checks.
package portretry

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portretry: invalid port number")

// Entry holds retry state for a single port.
type Entry struct {
	Attempts  int
	LastRetry time.Time
	NextRetry time.Time
}

// Retrier tracks retry attempts per port.
type Retrier struct {
	mu       sync.Mutex
	entries  map[int]*Entry
	backoff  time.Duration
	maxRetry int
}

// New creates a Retrier with the given backoff duration and max retry count.
func New(backoff time.Duration, maxRetry int) *Retrier {
	if backoff <= 0 {
		backoff = time.Second
	}
	if maxRetry <= 0 {
		maxRetry = 3
	}
	return &Retrier{
		entries:  make(map[int]*Entry),
		backoff:  backoff,
		maxRetry: maxRetry,
	}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Record increments the retry count for a port and schedules the next retry.
func (r *Retrier) Record(port int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	e, ok := r.entries[port]
	if !ok {
		e = &Entry{}
		r.entries[port] = e
	}
	e.Attempts++
	e.LastRetry = now
	e.NextRetry = now.Add(time.Duration(e.Attempts) * r.backoff)
	return nil
}

// Get returns the retry entry for a port, or nil if none exists.
func (r *Retrier) Get(port int) *Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	e := r.entries[port]
	if e == nil {
		return nil
	}
	copy := *e
	return &copy
}

// Exceeded reports whether the port has hit the max retry limit.
func (r *Retrier) Exceeded(port int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[port]
	return ok && e.Attempts >= r.maxRetry
}

// Reset clears retry state for a port.
func (r *Retrier) Reset(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}
