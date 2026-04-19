// Package portcount tracks how many times each port has been seen open across scans.
package portcount

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portcount: invalid port number")

// Counter tracks open-seen counts per port.
type Counter struct {
	mu     sync.RWMutex
	counts map[int]int
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{counts: make(map[int]int)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Inc increments the seen count for the given port.
func (c *Counter) Inc(port int) error {
	if !validPort(port) {
		return fmt.Errorf("%w: %d", ErrInvalidPort, port)
	}
	c.mu.Lock()
	c.counts[port]++
	c.mu.Unlock()
	return nil
}

// Get returns the current count for a port. Returns 0 if never seen.
func (c *Counter) Get(port int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counts[port]
}

// Reset clears the count for a port.
func (c *Counter) Reset(port int) {
	c.mu.Lock()
	delete(c.counts, port)
	c.mu.Unlock()
}

// All returns a copy of all port counts.
func (c *Counter) All() map[int]int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[int]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}
