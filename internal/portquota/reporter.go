package portquota

import (
	"fmt"
	"io"
	"sort"
)

// Reporter prints quota summaries to a writer.
type Reporter struct {
	quota *Quota
	out   io.Writer
}

// NewReporter returns a Reporter that reads from q and writes to out.
func NewReporter(q *Quota, out io.Writer) *Reporter {
	return &Reporter{quota: q, out: out}
}

// Print writes a sorted table of port quotas to the reporter's writer.
func (r *Reporter) Print() {
	r.quota.mu.RLock()
	ports := make([]int, 0, len(r.quota.entries))
	for p := range r.quota.entries {
		ports = append(ports, p)
	}
	r.quota.mu.RUnlock()

	if len(ports) == 0 {
		fmt.Fprintln(r.out, "no quota entries")
		return
	}

	sort.Ints(ports)
	fmt.Fprintf(r.out, "%-8s %-8s %-8s %s\n", "PORT", "HITS", "LIMIT", "EXCEEDED")
	for _, p := range ports {
		e, ok := r.quota.Get(p)
		if !ok {
			continue
		}
		exceeded := "no"
		if e.Exceeded() {
			exceeded = "YES"
		}
		fmt.Fprintf(r.out, "%-8d %-8d %-8d %s\n", p, e.Hits, e.Limit, exceeded)
	}
}
