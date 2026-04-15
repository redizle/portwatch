package monitor

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
)

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func newTestMonitor(t *testing.T, ports []int) *Monitor {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Ports = ports
	cfg.Interval = 50 * time.Millisecond

	log, _ := logger.New("stdout", "text")
	alerter := alert.New(nil)
	return New(cfg, log, alerter)
}

func TestMonitor_RunCancels(t *testing.T) {
	m := newTestMonitor(t, []int{9})
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- m.Run(ctx) }()

	time.Sleep(80 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("monitor did not stop after context cancel")
	}
}

func TestMonitor_DetectsOpenPort(t *testing.T) {
	port := freePort(t)
	ln, err := net.Listen("tcp", "127.0.0.1:"+itoa(port))
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	m := newTestMonitor(t, []int{port})
	m.scan()

	ps, ok := m.store.Get(port)
	if !ok {
		t.Fatal("expected port in store after scan")
	}
	if !ps.Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
