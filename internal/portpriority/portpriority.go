// Package portpriority assigns and manages priority levels for monitored ports.
package portpriority

import (
	"errors"
	"sync"
)

// Level represents a priority level for a port.
type Level int

const (
	Low    Level = 1
	Normal Level = 2
	High   Level = 3
	Critical Level = 4
)

func (l Level) String() string {
	switch l {
	case Low:
		return "low"
	case Normal:
		return "normal"
	case High:
		return "high"
	case Critical:
		return "critical"
	default:
		return "unknown"
	}
}

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("port must be between 1 and 65535")

// ErrInvalidLevel is returned when a priority level is unrecognised.
var ErrInvalidLevel = errors.New("invalid priority level")

// Registry stores priority levels keyed by port number.
type Registry struct {
	mu       sync.RWMutex
	entries  map[int]Level
	default_ Level
}

// New creates a new Registry with the given default level.
func New(defaultLevel Level) *Registry {
	return &Registry{
		entries:  make(map[int]Level),
		default_: defaultLevel,
	}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set assigns a priority level to a port.
func (r *Registry) Set(port int, level Level) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if level < Low || level > Critical {
		return ErrInvalidLevel
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[port] = level
	return nil
}

// Get returns the priority level for a port, falling back to the default.
func (r *Registry) Get(port int) Level {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if l, ok := r.entries[port]; ok {
		return l
	}
	return r.default_
}

// Remove deletes any explicit priority for a port.
func (r *Registry) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}

// All returns a copy of all explicitly set port priorities.
func (r *Registry) All() map[int]Level {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int]Level, len(r.entries))
	for k, v := range r.entries {
		out[k] = v
	}
	return out
}
