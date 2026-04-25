package portburst

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestRecord_InvalidPort(t *testing.T) {
	tr, _ := New(time.Second, 2)
	if err := tr.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := tr.Record(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestIsBursting_BelowThreshold(t *testing.T) {
	tr, _ := New(time.Second, 3)
	_ = tr.Record(8080)
	_ = tr.Record(8080)
	if tr.IsBursting(8080) {
		t.Fatal("should not be bursting below threshold")
	}
}

func TestIsBursting_AtThreshold(t *testing.T) {
	tr, _ := New(time.Second, 3)
	for i := 0; i < 3; i++ {
		_ = tr.Record(9090)
	}
	if !tr.IsBursting(9090) {
		t.Fatal("expected bursting at threshold")
	}
}

func TestIsBursting_MissingPort(t *testing.T) {
	tr, _ := New(time.Second, 2)
	if tr.IsBursting(1234) {
		t.Fatal("expected false for unknown port")
	}
}

func TestHitCount_WithinWindow(t *testing.T) {
	tr, _ := New(time.Second, 5)
	for i := 0; i < 4; i++ {
		_ = tr.Record(443)
	}
	if got := tr.HitCount(443); got != 4 {
		t.Fatalf("expected 4 hits, got %d", got)
	}
}

func TestHitCount_MissingPort(t *testing.T) {
	tr, _ := New(time.Second, 2)
	if got := tr.HitCount(9999); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsHits(t *testing.T) {
	tr, _ := New(time.Second, 2)
	_ = tr.Record(80)
	_ = tr.Record(80)
	tr.Reset(80)
	if tr.HitCount(80) != 0 {
		t.Fatal("expected zero hits after reset")
	}
	if tr.IsBursting(80) {
		t.Fatal("expected not bursting after reset")
	}
}

func TestIsBursting_ExcludesOldHits(t *testing.T) {
	tr, _ := New(50*time.Millisecond, 2)
	_ = tr.Record(5000)
	_ = tr.Record(5000)
	time.Sleep(60 * time.Millisecond)
	if tr.IsBursting(5000) {
		t.Fatal("old hits should be outside window")
	}
}
