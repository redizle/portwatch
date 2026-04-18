// Package portlabel assigns human-readable labels to ports for display and reporting.
package portlabel

import "sync"

// Label holds metadata for a port.
type Label struct {
	Port        int
	Name        string
	Description string
	Color       string // e.g. "red", "green", "yellow"
}

// Labeler manages port labels.
type Labeler struct {
	mu     sync.RWMutex
	labels map[int]Label
}

// New returns a new Labeler with optional seed labels.
func New(seed []Label) *Labeler {
	l := &Labeler{labels: make(map[int]Label)}
	for _, lb := range seed {
		if lb.Port > 0 && lb.Port <= 65535 {
			l.labels[lb.Port] = lb
		}
	}
	return l
}

// Set adds or replaces the label for a port.
func (l *Labeler) Set(lb Label) error {
	if lb.Port < 1 || lb.Port > 65535 {
		return ErrInvalidPort
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.labels[lb.Port] = lb
	return nil
}

// Get returns the label for a port, or a default if not set.
func (l *Labeler) Get(port int) (Label, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	lb, ok := l.labels[port]
	return lb, ok
}

// Remove deletes the label for a port.
func (l *Labeler) Remove(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.labels, port)
}

// All returns a copy of all labels.
func (l *Labeler) All() []Label {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Label, 0, len(l.labels))
	for _, lb := range l.labels {
		out = append(out, lb)
	}
	return out
}
