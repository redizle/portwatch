// Package portcensor provides redaction rules for sensitive ports,
// allowing certain ports to be hidden or masked in reports and logs.
package portcensor

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portcensor: port must be between 1 and 65535")

// ErrEmptyReason is returned when a blank redaction reason is provided.
var ErrEmptyReason = errors.New("portcensor: reason must not be empty")

// Entry holds the redaction metadata for a port.
type Entry struct {
	Port   int
	Reason string
}

// Censor manages a set of ports that should be redacted.
type Censor struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New creates an empty Censor.
func New() *Censor {
	return &Censor{entries: make(map[int]Entry)}
}

func validPort(p int) bool {
	return p >= 1 && p <= 65535
}

// Redact marks a port as censored with the given reason.
func (c *Censor) Redact(port int, reason string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if reason == "" {
		return ErrEmptyReason
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[port] = Entry{Port: port, Reason: reason}
	return nil
}

// Lift removes a port from the censored set.
func (c *Censor) Lift(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// IsCensored reports whether the given port is currently redacted.
func (c *Censor) IsCensored(port int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.entries[port]
	return ok
}

// Get returns the Entry for a port and whether it was found.
func (c *Censor) Get(port int) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[port]
	return e, ok
}

// All returns a snapshot of all censored entries.
func (c *Censor) All() []Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Entry, 0, len(c.entries))
	for _, e := range c.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of censored ports.
func (c *Censor) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
