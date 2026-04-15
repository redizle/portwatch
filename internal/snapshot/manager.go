package snapshot

import (
	"sync"
)

// Manager maintains the current and previous snapshots,
// allowing callers to rotate snapshots each scan cycle.
type Manager struct {
	mu      sync.Mutex
	current *Snapshot
	prev    *Snapshot
}

// NewManager creates a Manager with an empty initial snapshot.
func NewManager() *Manager {
	return &Manager{
		current: New(),
	}
}

// Current returns the active snapshot for writing.
func (m *Manager) Current() *Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.current
}

// Rotate promotes the current snapshot to previous and starts a fresh one.
// Returns the diff of changed ports between cycles.
func (m *Manager) Rotate() []PortState {
	m.mu.Lock()
	defer m.mu.Unlock()

	changed := m.current.Diff(m.prev)
	m.prev = m.current
	m.current = New()
	return changed
}

// Previous returns the snapshot from the last completed cycle.
func (m *Manager) Previous() *Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.prev
}
