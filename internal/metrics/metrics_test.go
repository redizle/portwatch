package metrics

import (
	"testing"
	"time"
)

func TestNew_SetsStartedAt(t *testing.T) {
	before := time.Now()
	m := New()
	after := time.Now()

	s := m.Snapshot()
	if s.StartedAt.Before(before) || s.StartedAt.After(after) {
		t.Errorf("StartedAt out of expected range: %v", s.StartedAt)
	}
}

func TestIncScans(t *testing.T) {
	m := New()
	m.IncScans()
	m.IncScans()
	if got := m.Snapshot().ScansTotal; got != 2 {
		t.Errorf("expected 2 scans, got %d", got)
	}
}

func TestIncOpen(t *testing.T) {
	m := New()
	m.IncOpen()
	if got := m.Snapshot().PortsOpen; got != 1 {
		t.Errorf("expected 1 open, got %d", got)
	}
}

func TestIncClosed(t *testing.T) {
	m := New()
	m.IncClosed()
	m.IncClosed()
	m.IncClosed()
	if got := m.Snapshot().PortsClosed; got != 3 {
		t.Errorf("expected 3 closed, got %d", got)
	}
}

func TestIncAlerts(t *testing.T) {
	m := New()
	m.IncAlerts()
	if got := m.Snapshot().AlertsSent; got != 1 {
		t.Errorf("expected 1 alert, got %d", got)
	}
}

func TestIncDropped(t *testing.T) {
	m := New()
	m.IncDropped()
	m.IncDropped()
	if got := m.Snapshot().FilterDropped; got != 2 {
		t.Errorf("expected 2 dropped, got %d", got)
	}
}

func TestUptime_Positive(t *testing.T) {
	m := New()
	time.Sleep(2 * time.Millisecond)
	if m.Uptime() <= 0 {
		t.Error("expected positive uptime")
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	m := New()
	s1 := m.Snapshot()
	m.IncScans()
	s2 := m.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Error("snapshot should be independent copy")
	}
}
