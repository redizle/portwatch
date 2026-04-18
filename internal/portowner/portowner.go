// Package portowner maps ports to logical owner names or team labels.
package portowner

import (
	"errors"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portowner: port must be between 1 and 65535")

// ErrEmptyOwner is returned when an empty owner string is provided.
var ErrEmptyOwner = errors.New("portowner: owner must not be empty")

// Registry holds port-to-owner mappings.
type Registry struct {
	mu     sync.RWMutex
	owners map[int]string
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{owners: make(map[int]string)}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Set assigns an owner to a port.
func (r *Registry) Set(port int, owner string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if owner == "" {
		return ErrEmptyOwner
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.owners[port] = owner
	return nil
}

// Get returns the owner for a port and whether it was found.
func (r *Registry) Get(port int) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	owner, ok := r.owners[port]
	return owner, ok
}

// Remove deletes the owner mapping for a port.
func (r *Registry) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.owners, port)
}

// All returns a copy of all port-to-owner mappings.
func (r *Registry) All() map[int]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[int]string, len(r.owners))
	for k, v := range r.owners {
		copy[k] = v
	}
	return copy
}
