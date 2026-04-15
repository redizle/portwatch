package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single port event recorded in history.
type Entry struct {
	Port      int       `json:"port"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// History manages an in-memory log of port events with optional file persistence.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	filePath string
	maxEntries int
}

// New creates a new History instance. filePath may be empty to disable persistence.
func New(filePath string, maxEntries int) *History {
	if maxEntries <= 0 {
		maxEntries = 1000
	}
	h := &History{
		filePath:   filePath,
		maxEntries: maxEntries,
	}
	if filePath != "" {
		_ = h.load()
	}
	return h
}

// Record appends a new entry and persists if a file path is configured.
func (h *History) Record(port int, status string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = append(h.entries, Entry{
		Port:      port,
		Status:    status,
		Timestamp: time.Now().UTC(),
	})

	if len(h.entries) > h.maxEntries {
		h.entries = h.entries[len(h.entries)-h.maxEntries:]
	}

	if h.filePath != "" {
		_ = h.persist()
	}
}

// GetAll returns a copy of all recorded entries.
func (h *History) GetAll() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// GetByPort returns entries filtered by port number.
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

func (h *History) persist() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
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
