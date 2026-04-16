// Package baseline tracks the expected "normal" set of open ports
// and flags deviations from that baseline.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single baselined port.
type Entry struct {
	Port      int       `json:"port"`
	AddedAt   time.Time `json:"added_at"`
	Note      string    `json:"note,omitempty"`
}

// Baseline holds the set of ports considered normal.
type Baseline struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an empty Baseline.
func New() *Baseline {
	return &Baseline{entries: make(map[int]Entry)}
}

// Add marks a port as baselined.
func (b *Baseline) Add(port int, note string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[port] = Entry{Port: port, AddedAt: time.Now(), Note: note}
}

// Remove removes a port from the baseline.
func (b *Baseline) Remove(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, port)
}

// Contains reports whether port is baselined.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[port]
	return ok
}

// Unexpected returns ports from active that are not in the baseline.
func (b *Baseline) Unexpected(active []int) []int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var out []int
	for _, p := range active {
		if _, ok := b.entries[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// All returns a copy of all baselined entries.
func (b *Baseline) All() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

// Save persists the baseline to a JSON file.
func (b *Baseline) Save(path string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b.entries)
}

// Load reads a baseline from a JSON file.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	entries := make(map[int]Entry)
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, err
	}
	return &Baseline{entries: entries}, nil
}
