package metrics

import (
	"testing"
	"time"
)

func TestSnapshot_FieldsMatchCounters(t *testing.T) {
	m := New()

	m.IncScans()
	m.IncScans()
	m.IncOpen()
	m.IncClosed()
	m.IncAlerts()
	m.IncErrors()
	m.IncStateChanges()
	m.IncNotifications()

	snap := m.Snapshot()

	if snap.Scans != 2 {
		t.Errorf("expected Scans=2, got %d", snap.Scans)
	}
	if snap.OpenPorts != 1 {
		t.Errorf("expected OpenPorts=1, got %d", snap.OpenPorts)
	}
	if snap.ClosedPorts != 1 {
		t.Errorf("expected ClosedPorts=1, got %d", snap.ClosedPorts)
	}
	if snap.Alerts != 1 {
		t.Errorf("expected Alerts=1, got %d", snap.Alerts)
	}
	if snap.Errors != 1 {
		t.Errorf("expected Errors=1, got %d", snap.Errors)
	}
	if snap.StateChanges != 1 {
		t.Errorf("expected StateChanges=1, got %d", snap.StateChanges)
	}
	if snap.Notifications != 1 {
		t.Errorf("expected Notifications=1, got %d", snap.Notifications)
	}
}

func TestSnapshot_StartedAtIsSet(t *testing.T) {
	before := time.Now()
	m := New()
	after := time.Now()

	snap := m.Snapshot()

	if snap.StartedAt.Before(before) || snap.StartedAt.After(after) {
		t.Errorf("StartedAt %v not within expected range [%v, %v]", snap.StartedAt, before, after)
	}
}

func TestSnapshot_UptimeNonEmpty(t *testing.T) {
	m := New()
	snap := m.Snapshot()

	if snap.Uptime == "" {
		t.Error("expected non-empty Uptime string")
	}
}

func TestSnapshot_IsImmutable(t *testing.T) {
	m := New()
	m.IncScans()

	snap1 := m.Snapshot()
	m.IncScans()
	snap2 := m.Snapshot()

	if snap1.Scans == snap2.Scans {
		t.Error("expected snapshots to differ after incrementing")
	}
	if snap1.Scans != 1 {
		t.Errorf("first snapshot should be unchanged, got Scans=%d", snap1.Scans)
	}
}
