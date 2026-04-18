// Package metrics provides lightweight runtime counters for the portwatch daemon.
//
// Usage:
//
//	m := metrics.New()
//	m.IncScans()
//	m.IncOpen()
//
//	// Print a summary to stdout
//	r := metrics.NewReporter(os.Stdout)
//	r.Print(m.Snapshot())
//
// All counter methods are safe for concurrent use. Snapshot returns an
// independent copy of the current state so callers can inspect values
// without holding any lock.
//
// Counter overview:
//
//	- Scans: total number of port-scan cycles completed.
//	- Open:  cumulative count of ports observed in the open state.
//	- Close: cumulative count of ports observed transitioning to closed.
//	- Err:   total errors encountered during scanning (e.g. permission
//	         denied, network unreachable).
package metrics
