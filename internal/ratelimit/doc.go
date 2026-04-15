// Package ratelimit implements a concurrency-safe, per-key rate limiter
// for use within portwatch to prevent alert and notification flooding.
//
// A Limiter is configured with a minimum interval between allowed events
// for any given key (typically a port identifier such as "port:8080").
// Callers invoke Allow to check whether an event should proceed; the
// limiter records the timestamp of each allowed event and suppresses
// subsequent calls that arrive before the interval has elapsed.
//
// Example usage:
//
//	l := ratelimit.New(30 * time.Second)
//	if l.Allow("port:443") {
//		// send alert
//	}
package ratelimit
