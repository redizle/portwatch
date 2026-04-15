package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port event.
type Entry struct {
	Port      int       `json:"port"`
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

// History stores a bounded log of port events with optional persistence.
type History struct {
	mu         sync.RWMutex
	entries    []Entry
	maxEntries int
	filePath   string
}

// New creates a History with the given capacity and optional file path.
func New(maxEntries int, filePath string) *History {
	h := &History{
		maxEntries: maxEntries,
		filePath:   filePath,
	}
	if filePath != "" {
		_ = h.load()
	}
	return h
}

// Record appends a new event entry for the given port.
func (h *History) Record(port int, event string, ts time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, Entry{Port: port, Event: event, Timestamp: ts})
	if len(h.entries) > h.maxEntries {
		h.entries = h.entries[len(h.entries)-h.maxEntries:]
	}
}

// GetByPort returns all entries for a specific port.
func (h *History) GetByPort(port int) []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var out []Entry
	for _, e := range h.entries {
		if e.Port == port {
			out = append(out, e)
		}
	}
	return out
}

// Recent returns up to n most recent entries across all ports.
func (h *History) Recent(n int) []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if n >= len(h.entries) {
		copy := make([]Entry, len(h.entries))
		copy = append(copy[:0], h.entries...)
		return copy
	}
	result := make([]Entry, n)
	copy(result, h.entries[len(h.entries)-n:])
	return result
}

// Persist writes the current history to disk.
func (h *History) Persist() error {
	if h.filePath == "" {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, err := json.Marshal(h.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(h.filePath, data, 0644)
}

func (h *History) load() error {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

// All returns a copy of all recorded entries.
func (h *History) All() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}
