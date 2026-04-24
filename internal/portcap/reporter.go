package portcap

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// Reporter prints peak capacity summaries to an io.Writer.
type Reporter struct {
	cap *Cap
	out io.Writer
}

// NewReporter returns a Reporter backed by c that writes to out.
func NewReporter(c *Cap, out io.Writer) *Reporter {
	return &Reporter{cap: c, out: out}
}

// Print writes a formatted table of all recorded peaks.
func (r *Reporter) Print() {
	all := r.cap.All()
	if len(all) == 0 {
		fmt.Fprintln(r.out, "no peak capacity data recorded")
		return
	}

	ports := make([]int, 0, len(all))
	for p := range all {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPEAK HITS\tRECORDED AT")
	for _, p := range ports {
		peak := all[p]
		fmt.Fprintf(tw, "%d\t%d\t%s\n", p, peak.Count,
			peak.RecordedAt.Format(time.RFC3339))
	}
	_ = tw.Flush()
}
