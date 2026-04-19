package portseq

import (
	"testing"
)

func TestRecord_AssignsSequence(t *testing.T) {
	s := New()
	e, created, err := s.Record(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected created=true")
	}
	if e.Sequence != 1 {
		t.Errorf("expected sequence 1, got %d", e.Sequence)
	}
}

func TestRecord_Idempotent(t *testing.T) {
	s := New()
	s.Record(443)
	e, created, err := s.Record(443)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Fatal("expected created=false on second record")
	}
	if e.Sequence != 1 {
		t.Errorf("expected sequence 1, got %d", e.Sequence)
	}
}

func TestRecord_MonotonicallyIncreasing(t *testing.T) {
	s := New()
	e1, _, _ := s.Record(80)
	e2, _, _ := s.Record(443)
	e3, _, _ := s.Record(22)
	if e1.Sequence >= e2.Sequence || e2.Sequence >= e3.Sequence {
		t.Errorf("sequences not increasing: %d %d %d", e1.Sequence, e2.Sequence, e3.Sequence)
	}
}

func TestRecord_InvalidPort(t *testing.T) {
	s := New()
	_, _, err := s.Record(0)
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, _, err = s.Record(70000)
	if err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Fatal("expected ok=false for unrecorded port")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	s := New()
	s.Record(80)
	s.Record(443)
	s.Reset()
	if s.Len() != 0 {
		t.Errorf("expected len 0 after reset, got %d", s.Len())
	}
	e, created, _ := s.Record(80)
	if !created || e.Sequence != 1 {
		t.Errorf("expected sequence reset to 1, got %d", e.Sequence)
	}
}

func TestLen_Counts(t *testing.T) {
	s := New()
	s.Record(80)
	s.Record(443)
	s.Record(80)
	if s.Len() != 2 {
		t.Errorf("expected len 2, got %d", s.Len())
	}
}
