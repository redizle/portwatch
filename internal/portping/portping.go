// Package portping provides latency measurement for open ports.
package portping

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single ping attempt.
type Result struct {
	Port    int
	Latency time.Duration
	Reachable bool
	Err     error
}

// Pinger measures TCP connect latency to localhost ports.
type Pinger struct {
	timeout time.Duration
}

// New returns a Pinger with the given connect timeout.
func New(timeout time.Duration) *Pinger {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Pinger{timeout: timeout}
}

// Ping attempts a TCP connection to the given port and returns latency.
func (p *Pinger) Ping(port int) Result {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Port: port, Latency: latency, Reachable: false, Err: err}
	}
	conn.Close()
	return Result{Port: port, Latency: latency, Reachable: true}
}

// PingAll pings each port in the slice and returns all results.
func (p *Pinger) PingAll(ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Ping(port))
	}
	return results
}
