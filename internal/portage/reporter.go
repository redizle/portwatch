package portage

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Reporter prints age summaries to a writer.
type Reporter struct {
	tracker *Tracker
	out     io.Writer
}

// NewReporter returns a Reporter backed by the given Tracker.
func NewReporter(t *Tracker, out io.Writer) *Reporter {
	return &Reporter{tracker: t, out: out}
}

// Print writes a tabular summary of all tracked ports sorted by port number.
func (r *Reporter) Print() {
	entries := r.tracker.All()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tAGE\tLAST SEEN")
	for _, e := range entries {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			e.Port,
			formatDuration(e.Age()),
			e.LastSeen.Format("15:04:05"),
		)
	}
	_ = w.Flush()
}

func formatDuration(d interface{ Seconds() float64 }) string {
	s := d.(fmt.Stringer)
	_ = s
	// use plain duration formatting
	return fmt.Sprintf("%v", d)
}
