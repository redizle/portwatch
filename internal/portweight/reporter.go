package portweight

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Reporter prints weight information to a writer.
type Reporter struct {
	weights *Weights
}

// NewReporter returns a Reporter backed by wt.
func NewReporter(wt *Weights) *Reporter {
	return &Reporter{weights: wt}
}

// Print writes a sorted table of port weights to w.
func (r *Reporter) Print(w io.Writer) {
	entries := r.weights.All()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no port weights configured")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tWEIGHT")
	for _, e := range entries {
		fmt.Fprintf(tw, "%d\t%d\n", e.Port, e.Weight)
	}
	_ = tw.Flush()
}
