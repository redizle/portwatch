// Package debounce delays repeated events for the same key until a quiet
// period has elapsed, reducing noise from rapidly flapping ports.
package debounce

import (
	"sync"
	"time"
)

// Debouncer holds pending events and fires them after a quiet window.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	timers  map[string]*time.Timer
	handler func(key string)
}

// New creates a Debouncer that waits window before calling handler.
func New(window time.Duration, handler func(key string)) *Debouncer {
	return &Debouncer{
		window:  window,
		timers:  make(map[string]*time.Timer),
		handler: handler,
	}
}

// Push schedules handler(key) after the quiet window. If Push is called
// again for the same key before the window expires the timer resets.
func (d *Debouncer) Push(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Reset(d.window)
		return
	}

	d.timers[key] = time.AfterFunc(d.window, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.handler(key)
	})
}

// Cancel removes a pending event for key without firing it.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys currently waiting to fire.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
