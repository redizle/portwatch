// Package portgroup provides named groupings of ports for easier monitoring and reporting.
package portgroup

import (
	"fmt"
	"sync"
)

// Group represents a named collection of ports.
type Group struct {
	Name  string
	Ports []int
}

// Registry holds named port groups.
type Registry struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string]*Group)}
}

// Add registers a named group of ports. Returns an error if any port is out of range.
func (r *Registry) Add(name string, ports []int) error {
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("portgroup: port %d out of range", p)
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := make([]int, len(ports))
	for i, p := range ports {
		copy[i] = p
	}
	r.groups[name] = &Group{Name: name, Ports: copy}
	return nil
}

// Remove deletes a group by name.
func (r *Registry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.groups, name)
}

// Get returns the Group for the given name, or false if not found.
func (r *Registry) Get(name string) (*Group, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.groups[name]
	return g, ok
}

// GroupsFor returns the names of all groups that contain the given port.
func (r *Registry) GroupsFor(port int) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name, g := range r.groups {
		for _, p := range g.Ports {
			if p == port {
				names = append(names, name)
				break
			}
		}
	}
	return names
}

// All returns a copy of all registered groups.
func (r *Registry) All() []*Group {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Group, 0, len(r.groups))
	for _, g := range r.groups {
		out = append(out, g)
	}
	return out
}
