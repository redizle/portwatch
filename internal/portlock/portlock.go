// Package portlock tracks ports that have been explicitly locked (ignored)
// by the operator so they are excluded from alerting and reporting.
package portlock

import (
	"fmt"
	"sync"
	"time"
)

// Entry describes a single locked port.
type Entry struct {
	Port      int
	Reason    string
	LockedAt  time.Time
	ExpiresAt *time.Time // nil means no expiry
}

// Locker manages the set of locked ports.
type Locker struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an empty Locker.
func New() *Locker {
	return &Locker{entries: make(map[int]Entry)}
}

// Lock adds or replaces a lock on the given port.
func (l *Locker) Lock(port int, reason string, ttl *time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portlock: invalid port %d", port)
	}
	e := Entry{
		Port:     port,
		Reason:   reason,
		LockedAt: time.Now(),
	}
	if ttl != nil {
		t := time.Now().Add(*ttl)
		e.ExpiresAt = &t
	}
	l.mu.Lock()
	l.entries[port] = e
	l.mu.Unlock()
	return nil
}

// Unlock removes the lock on the given port.
func (l *Locker) Unlock(port int) {
	l.mu.Lock()
	delete(l.entries, port)
	l.mu.Unlock()
}

// IsLocked reports whether the port is currently locked (and not expired).
func (l *Locker) IsLocked(port int) bool {
	l.mu.RLock()
	e, ok := l.entries[port]
	l.mu.RUnlock()
	if !ok {
		return false
	}
	if e.ExpiresAt != nil && time.Now().After(*e.ExpiresAt) {
		l.mu.Lock()
		delete(l.entries, port)
		l.mu.Unlock()
		return false
	}
	return true
}

// Active returns a snapshot of all currently locked (non-expired) entries.
func (l *Locker) Active() []Entry {
	l.mu.RLock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		if e.ExpiresAt == nil || time.Now().Before(*e.ExpiresAt) {
			out = append(out, e)
		}
	}
	l.mu.RUnlock()
	return out
}
