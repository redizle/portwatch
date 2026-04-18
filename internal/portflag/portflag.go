// Package portflag allows ports to be flagged with a short reason string
// for manual review or investigation purposes.
package portflag

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidPort is returned when the port number is out of range.
var ErrInvalidPort = errors.New("portflag: port must be between 1 and 65535")

// ErrEmptyReason is returned when an empty reason is provided.
var ErrEmptyReason = errors.New("portflag: reason must not be empty")

// Flag holds metadata about a flagged port.
type Flag struct {
	Port      int
	Reason    string
	FlaggedAt time.Time
}

// Flagger manages flagged ports.
type Flagger struct {
	mu    sync.RWMutex
	flags map[int]Flag
}

// New returns a new Flagger.
func New() *Flagger {
	return &Flagger{flags: make(map[int]Flag)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set flags a port with the given reason.
func (f *Flagger) Set(port int, reason string) error {
	if !validPort(port) {
		return ErrInvalidPort
	}
	if reason == "" {
		return ErrEmptyReason
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags[port] = Flag{Port: port, Reason: reason, FlaggedAt: time.Now()}
	return nil
}

// Get returns the flag for a port and whether it exists.
func (f *Flagger) Get(port int) (Flag, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	fl, ok := f.flags[port]
	return fl, ok
}

// IsFlagged reports whether a port is currently flagged.
func (f *Flagger) IsFlagged(port int) bool {
	_, ok := f.Get(port)
	return ok
}

// Unflag removes the flag from a port.
func (f *Flagger) Unflag(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.flags, port)
}

// All returns a copy of all flagged ports.
func (f *Flagger) All() []Flag {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]Flag, 0, len(f.flags))
	for _, fl := range f.flags {
		out = append(out, fl)
	}
	return out
}
