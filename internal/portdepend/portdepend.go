// Package portdepend tracks dependencies between ports,
// allowing users to declare that one port relies on another.
package portdepend

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("port must be between 1 and 65535")

// ErrSelfDependency is returned when a port is added as its own dependency.
var ErrSelfDependency = errors.New("port cannot depend on itself")

// Tracker stores port dependency relationships.
type Tracker struct {
	mu   sync.RWMutex
	deps map[int]map[int]struct{}
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{deps: make(map[int]map[int]struct{})}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Add records that port depends on dep.
func (t *Tracker) Add(port, dep int) error {
	if !validPort(port) || !validPort(dep) {
		return ErrInvalidPort
	}
	if port == dep {
		return ErrSelfDependency
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.deps[port] == nil {
		t.deps[port] = make(map[int]struct{})
	}
	t.deps[port][dep] = struct{}{}
	return nil
}

// Remove deletes a single dependency edge.
func (t *Tracker) Remove(port, dep int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.deps[port], dep)
}

// DepsOf returns all ports that port depends on.
func (t *Tracker) DepsOf(port int) []int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	set := t.deps[port]
	out := make([]int, 0, len(set))
	for d := range set {
		out = append(out, d)
	}
	return out
}

// Dependents returns all ports that depend on target.
func (t *Tracker) Dependents(target int) []int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []int
	for port, deps := range t.deps {
		if _, ok := deps[target]; ok {
			out = append(out, port)
		}
	}
	return out
}

// Clear removes all dependency records for port.
func (t *Tracker) Clear(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.deps, port)
}

// String returns a human-readable summary of dependencies for port.
func (t *Tracker) String(port int) string {
	deps := t.DepsOf(port)
	if len(deps) == 0 {
		return fmt.Sprintf("port %d has no dependencies", port)
	}
	return fmt.Sprintf("port %d depends on %v", port, deps)
}
