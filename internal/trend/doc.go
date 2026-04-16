// Package trend provides flap detection for monitored ports.
//
// It records open/close state-change events per port and exposes
// Churn and Flapping helpers so the daemon can suppress noisy alerts
// for ports that toggle rapidly within a sliding time window.
//
// Usage:
//
//	tr := trend.New(time.Minute, 100)
//	tr.Record(8080, true)
//	if tr.Flapping(8080, 5) {
//		// suppress or downgrade alert
//	}
package trend
