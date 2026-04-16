package uptime_test

import (
	"testing"
	"time"

	"portwatch/internal/uptime"
)

func TestMarkOpen_TracksPort(t *testing.T) {
	tr := uptime.New()
	tr.MarkOpen(8080)
	_, ok := tr.OpenedAt(8080)
	if !ok {
		t.Fatal("expected port 8080 to be tracked")
	}
}

func TestMarkOpen_Idempotent(t *testing.T) {
	tr := uptime.New()
	tr.MarkOpen(8080)
	first, _ := tr.OpenedAt(8080)
	time.Sleep(10 * time.Millisecond)
	tr.MarkOpen(8080)
	second, _ := tr.OpenedAt(8080)
	if !first.Equal(second) {
		t.Error("expected open time to remain unchanged on second MarkOpen")
	}
}

func TestMarkClosed_RemovesPort(t *testing.T) {
	tr := uptime.New()
	tr.MarkOpen(9090)
	tr.MarkClosed(9090)
	_, ok := tr.Uptime(9090)
	if ok {
		t.Error("expected port 9090 to be removed after MarkClosed")
	}
}

func TestUptime_ReturnsPositiveDuration(t *testing.T) {
	tr := uptime.New()
	tr.MarkOpen(443)
	time.Sleep(5 * time.Millisecond)
	d, ok := tr.Uptime(443)
	if !ok {
		t.Fatal("expected uptime to be tracked")
	}
	if d <= 0 {
		t.Errorf("expected positive duration, got %v", d)
	}
}

func TestUptime_MissingPort(t *testing.T) {
	tr := uptime.New()
	_, ok := tr.Uptime(1234)
	if ok {
		t.Error("expected false for untracked port")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := uptime.New()
	tr.MarkOpen(80)
	tr.MarkOpen(443)
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the copy should not affect tracker
	delete(all, 80)
	if _, ok := tr.OpenedAt(80); !ok {
		t.Error("tracker should not be affected by mutating All() result")
	}
}
