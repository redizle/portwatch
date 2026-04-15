package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestServer spins up a temporary TCP listener on a random port
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestCheckPort_Open(t *testing.T) {
	port, shutdown := startTestServer(t)
	defer shutdown()

	s := New("127.0.0.1", 2*time.Second)
	status := s.CheckPort(port)

	if !status.Open {
		t.Errorf("expected port %d to be open, got closed", port)
	}
	if status.Port != port {
		t.Errorf("expected port %d, got %d", port, status.Port)
	}
	if status.ScannedAt.IsZero() {
		t.Error("expected ScannedAt to be set")
	}
}

func TestCheckPort_Closed(t *testing.T) {
	s := New("127.0.0.1", 500*time.Millisecond)
	// port 1 is almost certainly closed/refused on localhost
	status := s.CheckPort(1)

	if status.Open {
		t.Error("expected port 1 to be closed")
	}
}

func TestScanPorts(t *testing.T) {
	port, shutdown := startTestServer(t)
	defer shutdown()

	s := New("127.0.0.1", 2*time.Second)
	results := s.ScanPorts([]int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
	if results[1].Open {
		t.Error("expected port 1 to be closed")
	}
}
