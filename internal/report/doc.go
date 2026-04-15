// Package report provides formatted output utilities for portwatch.
//
// It combines state and history data to produce human-readable
// summaries and activity logs, suitable for CLI output or
// periodic status dumps.
//
// Usage:
//
//	rep := report.New(historyInstance, stateInstance, os.Stdout)
//	rep.Summary()           // current port states
//	rep.RecentActivity(20)  // last 20 recorded events
package report
