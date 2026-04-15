// Package throttle provides a token-bucket style throttle for controlling
// how frequently port scan cycles are allowed to execute.
package throttle

import (
	"sync"
	"time"
)

// Throttle controls the rate of scan cycles using a minimum interval.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	lastRun  time.Time
	skipped  int64
}

// New creates a new Throttle with the given minimum interval between allowed calls.
func New(interval time.Duration) *Throttle {
	return &Throttle{
		interval: interval,
	}
}

// Allow returns true if enough time has passed since the last allowed call.
// If not enough time has passed, it increments the skipped counter and returns false.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if t.lastRun.IsZero() || now.Sub(t.lastRun) >= t.interval {
		t.lastRun = now
		return true
	}
	t.skipped++
	return false
}

// Skipped returns the total number of calls that were throttled.
func (t *Throttle) Skipped() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.skipped
}

// Reset clears the last run time, allowing the next call to Allow to pass immediately.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastRun = time.Time{}
}

// SetInterval updates the minimum interval between allowed calls.
func (t *Throttle) SetInterval(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.interval = d
}
