// Package throttle provides a simple interval-based throttle for controlling
// the frequency of repeated operations such as port scan cycles.
//
// It is safe for concurrent use. A Throttle tracks the last allowed call time
// and rejects calls that arrive before the configured interval has elapsed.
//
// Example usage:
//
//	th := throttle.New(5 * time.Second)
//	if th.Allow() {
//		// run scan
//	}
package throttle
