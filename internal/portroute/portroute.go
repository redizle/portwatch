// Package portroute maps ports to named routes or service endpoints,
// allowing operators to associate a port with a logical route identifier.
package portroute

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("port must be between 1 and 65535")

// ErrEmptyRoute is returned when an empty route string is provided.
var ErrEmptyRoute = errors.New("route must not be empty")

// Entry holds a route assignment for a port.
type Entry struct {
	Port  int
	Route string
}

// Router stores port-to-route mappings.
type Router struct {
	mu     sync.RWMutex
	routes map[int]string
}

// New returns an initialised Router.
func New() *Router {
	return &Router{routes: make(map[int]string)}
}

func validPort(p int) bool {
	return p >= 1 && p <= 65535
}

// Set assigns a route to the given port.
func (r *Router) Set(port int, route string) error {
	if !validPort(port) {
		return fmt.Errorf("%w: %d", ErrInvalidPort, port)
	}
	if route == "" {
		return ErrEmptyRoute
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[port] = route
	return nil
}

// Get returns the route for the given port and whether it was found.
func (r *Router) Get(port int) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.routes[port]
	return v, ok
}

// Remove deletes the route assignment for the given port.
func (r *Router) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, port)
}

// All returns a snapshot of all current route entries.
func (r *Router) All() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Entry, 0, len(r.routes))
	for p, route := range r.routes {
		out = append(out, Entry{Port: p, Route: route})
	}
	return out
}

// Len returns the number of route mappings.
func (r *Router) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.routes)
}
