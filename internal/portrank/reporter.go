package portrank

import (
	"fmt"
	"io"
	"sort"
)

// Reporter prints a ranked summary of ports to a writer.
type Reporter struct {
	ranker *Ranker
}

// NewReporter returns a Reporter backed by the given Ranker.
func NewReporter(r *Ranker) *Reporter {
	return &Reporter{ranker: r}
}

// Print writes a sorted-by-score table of all ranked ports to w.
func (rp *Reporter) Print(w io.Writer) {
	rp.ranker.mu.RLock()

	seen := make(map[int]struct{})
	for p := range rp.ranker.scores {
		seen[p] = struct{}{}
	}
	for p := range rp.ranker.override {
		seen[p] = struct{}{}
	}

	ports := make([]int, 0, len(seen))
	for p := range seen {
		ports = append(ports, p)
	}
	rp.ranker.mu.RUnlock()

	if len(ports) == 0 {
		fmt.Fprintln(w, "portrank: no entries")
		return
	}

	// Sort descending by effective score, then ascending by port.
	sort.Slice(ports, func(i, j int) bool {
		ei, _ := rp.ranker.Get(ports[i])
		ej, _ := rp.ranker.Get(ports[j])
		if ei.Score != ej.Score {
			return ei.Score > ej.Score
		}
		return ports[i] < ports[j]
	})

	fmt.Fprintf(w, "%-8s %-8s %s\n", "PORT", "SCORE", "OVERRIDE")
	for _, p := range ports {
		e, _ := rp.ranker.Get(p)
		ov := "-"
		if e.Override {
			ov = "yes"
		}
		fmt.Fprintf(w, "%-8d %-8d %s\n", p, e.Score, ov)
	}
}
