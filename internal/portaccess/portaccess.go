// Package portaccess tracks access control rules for monitored ports,
// allowing ports to be marked as allowed or denied with an optional reason.
package portaccess

import (
	"errors"
	"sync"
)

// Policy represents an access policy for a port.
type Policy int

const (
	PolicyAllow Policy = iota
	PolicyDeny
)

func (p Policy) String() string {
	switch p {
	case PolicyAllow:
		return "allow"
	case PolicyDeny:
		return "deny"
	default:
		return "unknown"
	}
}

// Entry holds the access policy and reason for a port.
type Entry struct {
	Policy Policy
	Reason string
}

// Tracker manages access control entries per port.
type Tracker struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[int]Entry),
	}
}

func validPort(port int) bool {
	return port >= 1 && port <= 65535
}

// Set assigns an access policy and reason to a port.
func (t *Tracker) Set(port int, policy Policy, reason string) error {
	if !validPort(port) {
		return errors.New("portaccess: port out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[port] = Entry{Policy: policy, Reason: reason}
	return nil
}

// Get returns the Entry for a port and whether it exists.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	return e, ok
}

// Remove deletes the access policy for a port.
func (t *Tracker) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// IsAllowed returns true if the port has an explicit allow policy.
// Ports with no entry default to allowed.
func (t *Tracker) IsAllowed(port int) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	if !ok {
		return true
	}
	return e.Policy == PolicyAllow
}

// Len returns the number of tracked entries.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}
