// Package portrelay forwards port events to registered downstream targets.
package portrelay

import (
	"errors"
	"fmt"
	"sync"
)

// Target receives forwarded port events.
type Target struct {
	Name    string
	Handler func(port int, event string) error
}

// Relay holds registered targets and dispatches events to them.
type Relay struct {
	mu      sync.RWMutex
	targets map[string]Target
}

// New returns an initialised Relay.
func New() *Relay {
	return &Relay{targets: make(map[string]Target)}
}

// Register adds a named target to the relay.
func (r *Relay) Register(t Target) error {
	if t.Name == "" {
		return errors.New("portrelay: target name must not be empty")
	}
	if t.Handler == nil {
		return errors.New("portrelay: target handler must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.targets[t.Name] = t
	return nil
}

// Unregister removes a target by name.
func (r *Relay) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.targets, name)
}

// Dispatch sends a port event to all registered targets.
// It collects and returns all errors encountered.
func (r *Relay) Dispatch(port int, event string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var errs []error
	for _, t := range r.targets {
		if err := t.Handler(port, event); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.Name, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("portrelay: %d target(s) failed: %v", len(errs), errs)
}

// Len returns the number of registered targets.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.targets)
}
