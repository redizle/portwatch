// Package portcap tracks the maximum observed hit count (capacity peak)
// for each monitored port within a rolling time window.
package portcap

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portcap: port must be between 1 and 65535")

// Peak holds the peak hit count and when it was recorded.
type Peak struct {
	Count     int
	RecordedAt time.Time
}

// Cap tracks peak hit counts per port.
type Cap struct {
	mu    sync.RWMutex
	peaks map[int]Peak
}

// New returns an initialised Cap.
func New() *Cap {
	return &Cap{peaks: make(map[int]Peak)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Observe records a hit count for the given port, updating the peak
// only when count exceeds the current maximum.
func (c *Cap) Observe(port, count int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.peaks[port]; !ok || count > existing.Count {
		c.peaks[port] = Peak{Count: count, RecordedAt: time.Now()}
	}
	return nil
}

// Get returns the current peak for a port and whether it exists.
func (c *Cap) Get(port int) (Peak, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.peaks[port]
	return p, ok
}

// Reset clears the peak record for a port.
func (c *Cap) Reset(port int) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.peaks, port)
	return nil
}

// All returns a copy of all recorded peaks keyed by port.
func (c *Cap) All() map[int]Peak {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[int]Peak, len(c.peaks))
	for k, v := range c.peaks {
		out[k] = v
	}
	return out
}
