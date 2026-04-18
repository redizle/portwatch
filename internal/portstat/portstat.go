// Package portstat tracks per-port scan statistics such as open/closed counts
// and last-seen timestamps.
package portstat

import (
	"fmt"
	"sync"
	"time"
)

// Stat holds cumulative statistics for a single port.
type Stat struct {
	Port       int
	OpenCount  int
	CloseCount int
	LastSeen   time.Time
	LastStatus string
}

// Tracker maintains statistics for observed ports.
type Tracker struct {
	mu    sync.RWMutex
	stats map[int]*Stat
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{stats: make(map[int]*Stat)}
}

// Record updates the statistic for port with the given status ("open" or "closed").
func (t *Tracker) Record(port int, status string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portstat: invalid port %d", port)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.stats[port]
	if !ok {
		s = &Stat{Port: port}
		t.stats[port] = s
	}
	s.LastSeen = time.Now()
	s.LastStatus = status
	switch status {
	case "open":
		s.OpenCount++
	case "closed":
		s.CloseCount++
	}
	return nil
}

// Get returns the Stat for port, or false if not found.
func (t *Tracker) Get(port int) (Stat, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.stats[port]
	if !ok {
		return Stat{}, false
	}
	return *s, true
}

// All returns a copy of all tracked stats.
func (t *Tracker) All() []Stat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Stat, 0, len(t.stats))
	for _, s := range t.stats {
		out = append(out, *s)
	}
	return out
}

// Reset clears all statistics.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stats = make(map[int]*Stat)
}
