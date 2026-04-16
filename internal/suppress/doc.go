// Package suppress provides per-port alert suppression for portwatch.
//
// During maintenance windows or known noisy periods, callers can suppress
// alerts for a specific port for a given duration. Suppression entries
// expire automatically and can also be lifted early.
//
// Example usage:
//
//	s := suppress.New()
//	s.Suppress(8080, "scheduled maintenance", 30*time.Minute)
//	if s.IsSuppressed(8080) {
//		// skip alert
//	}
package suppress
