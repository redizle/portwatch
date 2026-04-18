// Package portpause allows temporarily pausing monitoring for specific ports.
package portpause

import (
	"fmt"
	"sync"
	"time"
)

// Pauser tracks ports that are temporarily paused from monitoring.
type Pauser struct {
	mu      sync.RWMutex
	paused  map[int]time.Time // port -> resume time (zero = indefinite)
}

// New returns a new Pauser.
func New() *Pauser {
	return &Pauser{paused: make(map[int]time.Time)}
}

// Pause suspends monitoring for port for the given duration.
// Pass 0 for an indefinite pause.
func (p *Pauser) Pause(port int, d time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portpause: invalid port %d", port)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if d <= 0 {
		p.paused[port] = time.Time{}
	} else {
		p.paused[port] = time.Now().Add(d)
	}
	return nil
}

// Resume removes a pause for port immediately.
func (p *Pauser) Resume(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.paused, port)
}

// IsPaused reports whether port is currently paused.
func (p *Pauser) IsPaused(port int) bool {
	p.mu.RLock()
	resume, ok := p.paused[port]
	p.mu.RUnlock()
	if !ok {
		return false
	}
	if resume.IsZero() {
		return true
	}
	if time.Now().Before(resume) {
		return true
	}
	// expired — clean up
	p.mu.Lock()
	delete(p.paused, port)
	p.mu.Unlock()
	return false
}

// List returns all currently paused ports.
func (p *Pauser) List() []int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]int, 0, len(p.paused))
	for port, resume := range p.paused {
		if resume.IsZero() || time.Now().Before(resume) {
			out = append(out, port)
		}
	}
	return out
}
