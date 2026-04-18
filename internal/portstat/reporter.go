package portstat

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Reporter prints a formatted summary of port statistics.
type Reporter struct {
	tracker *Tracker
}

// NewReporter returns a Reporter backed by the given Tracker.
func NewReporter(t *Tracker) *Reporter {
	return &Reporter{tracker: t}
}

// Print writes a tabular summary of all tracked ports to w.
func (r *Reporter) Print(w io.Writer) error {
	stats := r.tracker.All()
	if len(stats) == 0 {
		_, err := fmt.Fprintln(w, "no port statistics recorded")
		return err
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Port < stats[j].Port
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tOPEN\tCLOSED\tLAST STATUS\tLAST SEEN")
	for _, s := range stats {
		fmt.Fprintf(tw, "%d\t%d\t%d\t%s\t%s\n",
			s.Port, s.OpenCount, s.CloseCount,
			s.LastStatus, s.LastSeen.Format("15:04:05"))
	}
	return tw.Flush()
}
