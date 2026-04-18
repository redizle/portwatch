// Package portexpiry tracks time-to-live (TTL) expiry for monitored ports.
// When a port's TTL elapses it is considered expired and can be evicted from
// active monitoring.
package portexpiry

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of the valid range.
var ErrInvalidPort = errors.New("portexpiry: port must be between 1 and 65535")

// Entry holds expiry metadata for a single port.
type Entry struct {
	Port      int
	ExpiresAt time.Time
}

// IsExpired reports whether the entry has passed its expiry time.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Registry manages TTL entries for ports.
type Registry struct {
	mu      sync.RWMutex
	entries map[int]Entry
	now     func() time.Time
}

// New creates a new Registry.
func New() *Registry {
	return &Registry{
		entries: make(map[int]Entry),
		now:     time.Now,
	}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set registers a TTL for the given port. Overwrites any existing entry.
func (r *Registry) Set(port int, ttl time.Duration) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[port] = Entry{Port: port, ExpiresAt: r.now().Add(ttl)}
	return nil
}

// Get returns the Entry for a port and whether it exists.
func (r *Registry) Get(port int) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[port]
	return e, ok
}

// Remove deletes the TTL entry for a port.
func (r *Registry) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}

// Expired returns all entries whose TTL has elapsed.
func (r *Registry) Expired() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Entry
	for _, e := range r.entries {
		if e.IsExpired() {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes all expired entries and returns them.
func (r *Registry) Evict() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var evicted []Entry
	for port, e := range r.entries {
		if time.Now().After(e.ExpiresAt) {
			evicted = append(evicted, e)
			delete(r.entries, port)
		}
	}
	return evicted
}
