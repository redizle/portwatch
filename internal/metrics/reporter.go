package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Reporter renders a Counters snapshot to an io.Writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a human-readable metrics summary to the underlying writer.
func (r *Reporter) Print(c Counters) error {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	uptime := time.Since(c.StartedAt).Round(time.Second)

	lines := []struct {
		label string
		value any
	}{
		{"Uptime", uptime},
		{"Scans total", c.ScansTotal},
		{"Ports open (events)", c.PortsOpen},
		{"Ports closed (events)", c.PortsClosed},
		{"Alerts sent", c.AlertsSent},
		{"Filter dropped", c.FilterDropped},
	}

	for _, l := range lines {
		if _, err := fmt.Fprintf(tw, "%s:\t%v\n", l.label, l.value); err != nil {
			return err
		}
	}
	return tw.Flush()
}
