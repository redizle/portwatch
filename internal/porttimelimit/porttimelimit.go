// Package porttimelimit enforces per-port time-based access windows.
// A port may be restricted to only be considered active during defined
// time ranges (e.g. business hours). Ports outside their window are
// treated as violations.
package porttimelimit

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Window defines a daily time range during which a port is permitted.
type Window struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// Entry holds the time window for a single port.
type Entry struct {
	Port   int
	Window Window
}

// Limiter stores per-port time windows.
type Limiter struct {
	mu      sync.RWMutex
	windows map[int]Window
}

// New returns an empty Limiter.
func New() *Limiter {
	return &Limiter{windows: make(map[int]Window)}
}

func validPort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("porttimelimit: invalid port %d", p)
	}
	return nil
}

// Set registers a time window for the given port.
func (l *Limiter) Set(port int, w Window) error {
	if err := validPort(port); err != nil {
		return err
	}
	if w.End <= w.Start {
		return errors.New("porttimelimit: end must be after start")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.windows[port] = w
	return nil
}

// Remove deletes the window for the given port.
func (l *Limiter) Remove(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.windows, port)
}

// Get returns the window for a port and whether it exists.
func (l *Limiter) Get(port int) (Window, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	w, ok := l.windows[port]
	return w, ok
}

// Allowed reports whether the given port is within its permitted window at t.
// If no window is registered for the port, Allowed returns true.
func (l *Limiter) Allowed(port int, t time.Time) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	w, ok := l.windows[port]
	if !ok {
		return true
	}
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	offset := t.Sub(midnight)
	return offset >= w.Start && offset < w.End
}

// Violations returns ports from active that fall outside their registered window at t.
func (l *Limiter) Violations(active []int, t time.Time) []int {
	var out []int
	for _, p := range active {
		if !l.Allowed(p, t) {
			out = append(out, p)
		}
	}
	return out
}
