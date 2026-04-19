// Package portnote stores timestamped notes attached to ports.
package portnote

import (
	"errors"
	"sync"
	"time"
)

// Note holds a single note entry for a port.
type Note struct {
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Store manages notes keyed by port number.
type Store struct {
	mu    sync.RWMutex
	notes map[int]*Note
}

// New returns an empty Store.
func New() *Store {
	return &Store{notes: make(map[int]*Note)}
}

func validPort(p int) bool { return p >= 1 && p <= 65535 }

// Set creates or replaces the note for port.
func (s *Store) Set(port int, text string) error {
	if !validPort(port) {
		return errors.New("portnote: port out of range")
	}
	if text == "" {
		return errors.New("portnote: text must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if n, ok := s.notes[port]; ok {
		n.Text = text
		n.UpdatedAt = now
	} else {
		s.notes[port] = &Note{Text: text, CreatedAt: now, UpdatedAt: now}
	}
	return nil
}

// Get returns the note for port and whether it exists.
func (s *Store) Get(port int) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notes[port]
	if !ok {
		return Note{}, false
	}
	return *n, true
}

// Remove deletes the note for port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.notes, port)
}

// Len returns the number of ports with notes.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.notes)
}
