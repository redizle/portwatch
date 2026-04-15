// Package metrics tracks runtime counters for portwatch daemon activity.
package metrics

import (
	"sync"
	"time"
)

// Counters holds cumulative stats collected during a daemon run.
type Counters struct {
	ScansTotal    int64
	PortsOpen     int64
	PortsClosed   int64
	Alertssent    int64
	FilterDropped int64
	StartedAt     time.Time
}

// Metrics is a thread-safe runtime metrics collector.
type Metrics struct {
	mu sync.RWMutex
	c  Counters
}

// New returns a new Metrics instance with StartedAt set to now.
func New() *Metrics {
	return &Metrics{
		c: Counters{StartedAt: time.Now()},
	}
}

// IncScans increments the total scan counter.
func (m *Metrics) IncScans() {
	m.mu.Lock()
	m.c.ScansTotal++
	m.mu.Unlock()
}

// IncOpen increments the open-port counter.
func (m *Metrics) IncOpen() {
	m.mu.Lock()
	m.c.PortsOpen++
	m.mu.Unlock()
}

// IncClosed increments the closed-port counter.
func (m *Metrics) IncClosed() {
	m.mu.Lock()
	m.c.PortsClosed++
	m.mu.Unlock()
}

// IncAlerts increments the alerts-sent counter.
func (m *Metrics) IncAlerts() {
	m.mu.Lock()
	m.c.AlertsSent++
	m.mu.Unlock()
}

// IncDropped increments the filter-dropped counter.
func (m *Metrics) IncDropped() {
	m.mu.Lock()
	m.c.FilterDropped++
	m.mu.Unlock()
}

// Snapshot returns a copy of the current counters.
func (m *Metrics) Snapshot() Counters {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.c
}

// Uptime returns the duration since the metrics were created.
func (m *Metrics) Uptime() time.Duration {
	return time.Since(m.c.StartedAt)
}
