// Package portscope assigns scope classifications to ports (internal, external, loopback, etc.)
package portscope

import (
	"errors"
	"sync"
)

// Scope represents a port's network scope classification.
type Scope string

const (
	ScopeInternal Scope = "internal"
	ScopeExternal Scope = "external"
	ScopeLoopback Scope = "loopback"
	ScopeUnknown  Scope = "unknown"
)

// Entry holds the scope for a port.
type Entry struct {
	Port  int
	Scope Scope
}

// Classifier manages port scope assignments.
type Classifier struct {
	mu      sync.RWMutex
	scopes  map[int]Scope
	default_ Scope
}

// New creates a new Classifier with the given default scope.
func New(defaultScope Scope) *Classifier {
	return &Classifier{
		scopes:   make(map[int]Scope),
		default_: defaultScope,
	}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set assigns a scope to a port.
func (c *Classifier) Set(port int, scope Scope) error {
	if !validPort(port) {
		return errors.New("portscope: port out of range")
	}
	if scope == "" {
		return errors.New("portscope: scope must not be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scopes[port] = scope
	return nil
}

// Get returns the scope for a port, falling back to the default.
func (c *Classifier) Get(port int) Scope {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if s, ok := c.scopes[port]; ok {
		return s
	}
	return c.default_
}

// Remove deletes the scope override for a port.
func (c *Classifier) Remove(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.scopes, port)
}

// All returns all explicitly set port scope entries.
func (c *Classifier) All() []Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Entry, 0, len(c.scopes))
	for p, s := range c.scopes {
		out = append(out, Entry{Port: p, Scope: s})
	}
	return out
}
