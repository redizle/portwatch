package portaudit

import (
	"testing"
)

func TestRecord_AddsEntry(t *testing.T) {
	l := New(0)
	l.Record(8080, ActionOpened, "test")
	if len(l.All()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.All()))
	}
}

func TestRecord_MaxSize_Evicts(t *testing.T) {
	l := New(3)
	for i := 0; i < 5; i++ {
		l.Record(i, ActionOpened, "")
	}
	if len(l.All()) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(l.All()))
	}
}

func TestForPort_Filters(t *testing.T) {
	l := New(0)
	l.Record(80, ActionOpened, "")
	l.Record(443, ActionOpened, "")
	l.Record(80, ActionClosed, "")
	res := l.ForPort(80)
	if len(res) != 2 {
		t.Fatalf("expected 2 entries for port 80, got %d", len(res))
	}
}

func TestForPort_Missing(t *testing.T) {
	l := New(0)
	if len(l.ForPort(9999)) != 0 {
		t.Fatal("expected empty result for unknown port")
	}
}

func TestClear_RemovesAll(t *testing.T) {
	l := New(0)
	l.Record(22, ActionAlerted, "")
	l.Clear()
	if len(l.All()) != 0 {
		t.Fatal("expected empty log after clear")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New(0)
	l.Record(80, ActionOpened, "")
	a := l.All()
	a[0].Port = 9999
	if l.All()[0].Port == 9999 {
		t.Fatal("All() should return a copy")
	}
}

func TestEntry_String(t *testing.T) {
	l := New(0)
	l.Record(8080, ActionSuppressed, "quiet hours")
	s := l.All()[0].String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
