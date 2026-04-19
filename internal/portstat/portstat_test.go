package portstat

import (
	"testing"
)

func TestRecord_Open(t *testing.T) {
	tr := New()
	if err := tr.Record(80, "open"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected stat for port 80")
	}
	if s.OpenCount != 1 {
		t.Errorf("expected OpenCount 1, got %d", s.OpenCount)
	}
	if s.LastStatus != "open" {
		t.Errorf("expected LastStatus open, got %s", s.LastStatus)
	}
}

func TestRecord_Closed(t *testing.T) {
	tr := New()
	_ = tr.Record(443, "closed")
	s, _ := tr.Get(443)
	if s.CloseCount != 1 {
		t.Errorf("expected CloseCount 1, got %d", s.CloseCount)
	}
}

func TestRecord_Accumulates(t *testing.T) {
	tr := New()
	for i := 0; i < 3; i++ {
		_ = tr.Record(8080, "open")
	}
	_ = tr.Record(8080, "closed")
	s, _ := tr.Get(8080)
	if s.OpenCount != 3 {
		t.Errorf("expected 3 opens, got %d", s.OpenCount)
	}
	if s.CloseCount != 1 {
		t.Errorf("expected 1 close, got %d", s.CloseCount)
	}
}

func TestRecord_InvalidPort(t *testing.T) {
	tr := New()
	if err := tr.Record(0, "open"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := tr.Record(70000, "open"); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestRecord_InvalidStatus(t *testing.T) {
	tr := New()
	if err := tr.Record(80, "unknown"); err == nil {
		t.Error("expected error for invalid status 'unknown'")
	}
	if err := tr.Record(80, ""); err == nil {
		t.Error("expected error for empty status")
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New()
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := New()
	_ = tr.Record(22, "open")
	_ = tr.Record(80, "open")
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 stats, got %d", len(all))
	}
}

func TestReset_ClearsStats(t *testing.T) {
	tr := New()
	_ = tr.Record(22, "open")
	tr.Reset()
	if len(tr.All()) != 0 {
		t.Error("expected empty stats after reset")
	}
}
