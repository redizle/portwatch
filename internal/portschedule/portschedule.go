// Package portschedule allows scheduling port scans at configurable intervals
// per port or port range, overriding the global scan interval.
package portschedule

import (
	"fmt"
	"sync"
	"time"
)

// Entry defines a custom scan schedule for a specific port.
type Entry struct {
	Port     int
	Interval time.Duration
	LastScan time.Time
}

// Schedule holds per-port scan intervals.
type Schedule struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an empty Schedule.
func New() *Schedule {
	return &Schedule{
		entries: make(map[int]*Entry),
	}
}

// Set registers a custom interval for the given port.
func (s *Schedule) Set(port int, interval time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portschedule: invalid port %d", port)
	}
	if interval <= 0 {
		return fmt.Errorf("portschedule: interval must be positive")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = &Entry{Port: port, Interval: interval}
	return nil
}

// Remove deletes the custom schedule for a port.
func (s *Schedule) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Due reports whether the given port is due for a scan based on its schedule.
// If no custom schedule exists, Due returns false.
func (s *Schedule) Due(port int, now time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[port]
	if !ok {
		return false
	}
	return now.Sub(e.LastScan) >= e.Interval
}

// MarkScanned updates the last scan time for the given port.
func (s *Schedule) MarkScanned(port int, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[port]; ok {
		e.LastScan = now
	}
}

// All returns a copy of all scheduled entries.
func (s *Schedule) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}
