// Package portstat tracks per-port scan statistics including open and closed
// counts and the last observed status for each port.
//
// Usage:
//
//	tr := portstat.New()
//	tr.Record(80, "open")
//	tr.Record(80, "closed")
//	s, ok := tr.Get(80)
//
// A Reporter can print a formatted table of all collected statistics to any
// io.Writer for quick inspection during a monitoring session.
package portstat
