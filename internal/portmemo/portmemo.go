// Package portmemo stores short-lived memos attached to ports.
package portmemo

import (
	"errors"
	"sync"
	"time"
)

// Memo holds a text note and an optional expiry.
type Memo struct {
	Text      string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

// Expired reports whether the memo has passed its expiry time.
func (m Memo) Expired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*m.ExpiresAt)
}

// Store manages memos keyed by port number.
type Store struct {
	mu    sync.RWMutex
	items map[int]Memo
}

// New returns an empty Store.
func New() *Store {
	return &Store{items: make(map[int]Memo)}
}

// Set attaches a memo to a port. ttl of 0 means no expiry.
func (s *Store) Set(port int, text string, ttl time.Duration) error {
	if port < 1 || port > 65535 {
		return errors.New("portmemo: port out of range")
	}
	if text == "" {
		return errors.New("portmemo: text must not be empty")
	}
	m := Memo{Text: text, CreatedAt: time.Now()}
	if ttl > 0 {
		t := time.Now().Add(ttl)
		m.ExpiresAt = &t
	}
	s.mu.Lock()
	s.items[port] = m
	s.mu.Unlock()
	return nil
}

// Get returns the memo for a port. Returns false if missing or expired.
func (s *Store) Get(port int) (Memo, bool) {
	s.mu.RLock()
	m, ok := s.items[port]
	s.mu.RUnlock()
	if !ok || m.Expired() {
		return Memo{}, false
	}
	return m, true
}

// Remove deletes the memo for a port.
func (s *Store) Remove(port int) {
	s.mu.Lock()
	delete(s.items, port)
	s.mu.Unlock()
}

// Purge removes all expired memos.
func (s *Store) Purge() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for p, m := range s.items {
		if m.Expired() {
			delete(s.items, p)
			n++
		}
	}
	return n
}
