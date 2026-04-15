package report

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/state"
)

// Reporter generates human-readable summaries of port activity.
type Reporter struct {
	history *history.History
	state   *state.State
	out     io.Writer
}

// New creates a new Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(h *history.History, s *state.State, w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{history: h, state: s, out: w}
}

// Summary prints a snapshot of currently tracked port states.
func (r *Reporter) Summary() {
	ports := r.state.All()
	if len(ports) == 0 {
		fmt.Fprintln(r.out, "No ports currently tracked.")
		return
	}

	keys := make([]int, 0, len(ports))
	for p := range ports {
		keys = append(keys, p)
	}
	sort.Ints(keys)

	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tSTATUS\tLAST SEEN")
	fmt.Fprintln(tw, "----\t------\t---------")
	for _, p := range keys {
		entry := ports[p]
		status := "closed"
		if entry.Open {
			status = "open"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\n", p, status, entry.LastSeen.Format(time.RFC3339))
	}
	tw.Flush()
}

// RecentActivity prints the last n history entries across all ports.
func (r *Reporter) RecentActivity(n int) {
	entries := r.history.Recent(n)
	if len(entries) == 0 {
		fmt.Fprintln(r.out, "No recent activity recorded.")
		return
	}

	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tPORT\tEVENT")
	fmt.Fprintln(tw, "----\t----\t-----")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\n", e.Timestamp.Format(time.RFC3339), e.Port, e.Event)
	}
	tw.Flush()
}
