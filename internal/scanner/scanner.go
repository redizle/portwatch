package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortStatus represents the state of a scanned port
type PortStatus struct {
	Port     int
	Open     bool
	Latency  time.Duration
	ScannedAt time.Time
}

// Scanner holds configuration for port scanning
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner with the given host and timeout
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: timeout,
	}
}

// CheckPort checks whether a single TCP port is open on the scanner's host
func (s *Scanner) CheckPort(port int) PortStatus {
	addr := fmt.Sprintf("%s:%d", s.Host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", addr, s.Timeout)
	latency := time.Since(start)

	status := PortStatus{
		Port:      port,
		Latency:   latency,
		ScannedAt: time.Now(),
	}

	if err != nil {
		status.Open = false
		return status
	}

	conn.Close()
	status.Open = true
	return status
}

// ScanPorts scans a slice of ports and returns their statuses
func (s *Scanner) ScanPorts(ports []int) []PortStatus {
	results := make([]PortStatus, 0, len(ports))
	for _, port := range ports {
		results = append(results, s.CheckPort(port))
	}
	return results
}
