package portping_test

import (
	"net"
	"testing"
	"time"

	"portwatch/internal/portping"
)

func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestPing_ReachablePort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	p := portping.New(time.Second)
	r := p.Ping(port)

	if !r.Reachable {
		t.Fatalf("expected reachable, got err: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Error("expected positive latency")
	}
	if r.Port != port {
		t.Errorf("port mismatch: got %d want %d", r.Port, port)
	}
}

func TestPing_ClosedPort(t *testing.T) {
	p := portping.New(200 * time.Millisecond)
	r := p.Ping(1)

	if r.Reachable {
		t.Error("expected unreachable")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestPingAll_MixedPorts(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	p := portping.New(500 * time.Millisecond)
	results := p.PingAll([]int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Reachable {
		t.Error("first port should be reachable")
	}
	if results[1].Reachable {
		t.Error("second port should not be reachable")
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	p := portping.New(0)
	if p == nil {
		t.Fatal("expected non-nil pinger")
	}
}
