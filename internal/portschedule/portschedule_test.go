package portschedule

import (
	"testing"
	"time"
)

func TestSet_ValidPort(t *testing.T) {
	s := New()
	if err := s.Set(8080, time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.All()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.All()))
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, time.Second); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Set(70000, time.Second); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestSet_ZeroInterval(t *testing.T) {
	s := New()
	if err := s.Set(80, 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(443, time.Minute)
	s.Remove(443)
	if len(s.All()) != 0 {
		t.Fatal("expected empty schedule after remove")
	}
}

func TestDue_NoEntry(t *testing.T) {
	s := New()
	if s.Due(9000, time.Now()) {
		t.Fatal("expected false for unscheduled port")
	}
}

func TestDue_NotYet(t *testing.T) {
	s := New()
	_ = s.Set(22, 10*time.Minute)
	now := time.Now()
	s.MarkScanned(22, now)
	if s.Due(22, now.Add(time.Second)) {
		t.Fatal("expected not due so soon")
	}
}

func TestDue_AfterInterval(t *testing.T) {
	s := New()
	_ = s.Set(22, 5*time.Second)
	past := time.Now().Add(-10 * time.Second)
	s.MarkScanned(22, past)
	if !s.Due(22, time.Now()) {
		t.Fatal("expected port to be due")
	}
}

func TestMarkScanned_UpdatesTime(t *testing.T) {
	s := New()
	_ = s.Set(3306, time.Minute)
	now := time.Now()
	s.MarkScanned(3306, now)
	entries := s.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if !entries[0].LastScan.Equal(now) {
		t.Errorf("LastScan not updated correctly")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, time.Second)
	_ = s.Set(443, time.Second)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}
