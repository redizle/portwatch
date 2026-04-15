// Package ratelimit provides a simple token-bucket style rate limiter
// used to throttle repeated alerts or scan events for the same port.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-key event rates and allows callers to check whether
// an event should be allowed through based on a minimum interval.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
}

// New creates a Limiter that allows at most one event per key per interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{
		interval: interval,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the configured
// interval, and records the current time for that key. Returns false if
// the key was seen too recently.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.interval {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for a key, allowing the next call to
// Allow for that key to pass immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all recorded keys, resetting the limiter entirely.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
