// Package trend tracks port open/close frequency over time.
package trend

import (
	"sync"
	"time"
)

// Entry records a single state change event for a port.
type Entry struct {
	Port      int
	Open      bool
	Timestamp time.Time
}

// Trend holds recent events per port and computes churn metrics.
type Trend struct {
	mu      sync.Mutex
	events  map[int][]Entry
	window  time.Duration
	maxPer  int
}

// New creates a Trend tracker. window defines how far back to look;
// maxPer caps stored events per port.
func New(window time.Duration, maxPer int) *Trend {
	if maxPer <= 0 {
		maxPer = 100
	}
	return &Trend{
		events: make(map[int][]Entry),
		window: window,
		maxPer: maxPer,
	}
}

// Record appends a state-change event for the given port.
func (t *Trend) Record(port int, open bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e := Entry{Port: port, Open: open, Timestamp: time.Now()}
	t.events[port] = append(t.events[port], e)
	if len(t.events[port]) > t.maxPer {
		t.events[port] = t.events[port][1:]
	}
}

// Churn returns the number of state changes for port within the window.
func (t *Trend) Churn(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := time.Now().Add(-t.window)
	count := 0
	for _, e := range t.events[port] {
		if e.Timestamp.After(cutoff) {
			count++
		}
	}
	return count
}

// Flapping returns true if the port has churned more than threshold times
// within the configured window.
func (t *Trend) Flapping(port int, threshold int) bool {
	return t.Churn(port) >= threshold
}

// Reset clears all recorded events for a port.
func (t *Trend) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.events, port)
}
