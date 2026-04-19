// Package portcooldown tracks per-port cooldown periods to suppress
// repeated alerts within a configurable window.
package portcooldown

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portcooldown: port must be between 1 and 65535")

// Entry holds cooldown state for a single port.
type Entry struct {
	Until     time.Time
	Triggered int
}

// Cooldown manages per-port cooldown windows.
type Cooldown struct {
	mu       sync.Mutex
	entries  map[int]Entry
	default_ time.Duration
}

// New creates a Cooldown with the given default window duration.
func New(defaultWindow time.Duration) *Cooldown {
	return &Cooldown{
		entries:  make(map[int]Entry),
		default_: defaultWindow,
	}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// IsCooling returns true if the port is currently within a cooldown window.
func (c *Cooldown) IsCooling(port int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[port]
	if !ok {
		return false
	}
	return time.Now().Before(e.Until)
}

// Trigger marks the port as cooling down using the default window.
// Returns ErrInvalidPort if port is out of range.
func (c *Cooldown) Trigger(port int) error {
	return c.TriggerFor(port, c.default_)
}

// TriggerFor marks the port as cooling down for the given duration.
func (c *Cooldown) TriggerFor(port int, d time.Duration) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.entries[port]
	e.Until = time.Now().Add(d)
	e.Triggered++
	c.entries[port] = e
	return nil
}

// Reset clears the cooldown for a port.
func (c *Cooldown) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// Get returns the Entry for a port and whether it exists.
func (c *Cooldown) Get(port int) (Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[port]
	return e, ok
}
