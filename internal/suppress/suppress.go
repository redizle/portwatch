// Package suppress provides a mechanism to temporarily silence alerts
// for specific ports, useful during maintenance windows or known noisy periods.
package suppress

import (
	"sync"
	"time"
)

// Entry holds suppression details for a port.
type Entry struct {
	Port      int
	Reason    string
	ExpiresAt time.Time
}

// Suppressor manages per-port alert suppression.
type Suppressor struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{entries: make(map[int]Entry)}
}

// Suppress silences alerts for port for the given duration.
func (s *Suppressor) Suppress(port int, reason string, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = Entry{
		Port:      port,
		Reason:    reason,
		ExpiresAt: time.Now().Add(d),
	}
}

// IsSuppressed reports whether alerts for port are currently suppressed.
func (s *Suppressor) IsSuppressed(port int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[port]
	if !ok {
		return false
	}
	if time.Now().After(e.ExpiresAt) {
		return false
	}
	return true
}

// Lift removes suppression for port immediately.
func (s *Suppressor) Lift(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Active returns all currently active suppression entries.
func (s *Suppressor) Active() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		if now.Before(e.ExpiresAt) {
			out = append(out, e)
		}
	}
	return out
}
