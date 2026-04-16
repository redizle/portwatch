// Package digest produces periodic scan summaries for reporting and alerting.
package digest

import (
	"fmt"
	"strings"
	"time"
)

// Entry holds a single port event in the digest.
type Entry struct {
	Port   int
	Status string
	Name   string
	At     time.Time
}

// Digest aggregates entries over a window.
type Digest struct {
	window  time.Duration
	entries []Entry
}

// New returns a Digest with the given aggregation window.
func New(window time.Duration) *Digest {
	return &Digest{window: window}
}

// Add appends an entry to the digest.
func (d *Digest) Add(e Entry) {
	if e.At.IsZero() {
		e.At = time.Now()
	}
	d.entries = append(d.entries, e)
}

// Flush returns all entries within the window and clears the buffer.
func (d *Digest) Flush() []Entry {
	cutoff := time.Now().Add(-d.window)
	var out []Entry
	for _, e := range d.entries {
		if e.At.After(cutoff) {
			out = append(out, e)
		}
	}
	d.entries = nil
	return out
}

// Summary returns a human-readable summary string.
func (d *Digest) Summary() string {
	if len(d.entries) == 0 {
		return "no activity"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d event(s):\n", len(d.entries))
	for _, e := range d.entries {
		fmt.Fprintf(&sb, "  port %d (%s): %s at %s\n", e.Port, e.Name, e.Status, e.At.Format(time.RFC3339))
	}
	return sb.String()
}

// Len returns the number of buffered entries.
func (d *Digest) Len() int { return len(d.entries) }
