// Package portquota provides per-port hit counting against configurable
// thresholds. It is useful for detecting when a monitored port receives
// an unexpectedly high volume of connection attempts within a scan cycle.
//
// Basic usage:
//
//	q := portquota.New()
//	_ = q.Set(8080, 100)   // allow up to 100 hits
//	_ = q.Inc(8080)        // record a hit
//	e, _ := q.Get(8080)
//	if e.Exceeded() {
//	    // trigger alert
//	}
//
// Reporters can be used to print a summary table:
//
//	r := portquota.NewReporter(q, os.Stdout)
//	r.Print()
package portquota
