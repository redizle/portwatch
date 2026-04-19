// Package portcooldown provides per-port cooldown tracking for portwatch.
//
// It allows callers to suppress repeated notifications for a port within
// a configurable time window. Each port tracks how many times it has been
// triggered and when its current cooldown expires.
//
// Example usage:
//
//	cd := portcooldown.New(30 * time.Second)
//	if !cd.IsCooling(8080) {
//		// send alert
//		_ = cd.Trigger(8080)
//	}
package portcooldown
