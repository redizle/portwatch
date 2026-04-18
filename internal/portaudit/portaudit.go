// Package portaudit records a timestamped audit trail of port-related
// actions (open, close, suppressed, alerted) for later review.
package portaudit

import (
	"fmt"
	"sync"
	"time"
)

// Action describes what happened to a port.
type Action string

const (
	ActionOpened     Action = "opened"
	ActionClosed     Action = "closed"
	ActionSuppressed Action = "suppressed"
	ActionAlerted    Action = "alerted"
)

// Entry is a single audit record.
type Entry struct {
	Port      int
	Action    Action
	Note      string
	Timestamp time.Time
}

func (e Entry) String() string {
	return fmt.Sprintf("%s port=%d action=%s note=%q", e.Timestamp.Format(time.RFC3339), e.Port, e.Action, e.Note)
}

// Log holds audit entries in memory up to a configurable cap.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates an audit Log capped at maxSize entries (0 = unlimited).
func New(maxSize int) *Log {
	return &Log{maxSize: maxSize}
}

// Record appends an entry to the log.
func (l *Log) Record(port int, action Action, note string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{Port: port, Action: action, Note: note, Timestamp: time.Now()}
	l.entries = append(l.entries, e)
	if l.maxSize > 0 && len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// All returns a copy of all current entries.
func (l *Log) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// ForPort returns entries filtered by port number.
func (l *Log) ForPort(port int) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Port == port {
			out = append(out, e)
		}
	}
	return out
}

// Clear removes all entries.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}
