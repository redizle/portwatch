// Package debounce provides a key-based debouncer that delays event
// delivery until a configurable quiet window has elapsed with no
// further pushes for the same key.
//
// Typical use-case: suppress repeated open/close notifications for a
// port that is flapping, only forwarding the event once the port state
// has been stable for the configured duration.
//
// Usage:
//
//	d := debounce.New(500*time.Millisecond, func(key string) {
//	    fmt.Println("stable:", key)
//	})
//	d.Push("8080")
package debounce
