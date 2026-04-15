package metrics

import "time"

// Snapshot holds a point-in-time copy of all metric counters.
type Snapshot struct {
	StartedAt  time.Time `json:"started_at"`
	Uptime     string    `json:"uptime"`
	Scans      uint64    `json:"scans"`
	OpenPorts  uint64    `json:"open_ports"`
	ClosedPorts uint64   `json:"closed_ports"`
	Alerts     uint64    `json:"alerts"`
	Errors     uint64    `json:"errors"`
	StateChanges uint64  `json:"state_changes"`
	Notifications uint64 `json:"notifications"`
}

// Snapshot returns an immutable copy of the current metrics.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Snapshot{
		StartedAt:     m.startedAt,
		Uptime:        time.Since(m.startedAt).Round(time.Second).String(),
		Scans:         m.scans,
		OpenPorts:     m.open,
		ClosedPorts:   m.closed,
		Alerts:        m.alerts,
		Errors:        m.errors,
		StateChanges:  m.stateChanges,
		Notifications: m.notifications,
	}
}
