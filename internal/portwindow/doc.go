// Package portwindow provides rolling time-window tracking for port observations.
//
// A Tracker records the first and last observation time for each port within
// a configurable window duration. When an observation arrives after the window
// has expired, the window resets from that point.
//
// Example usage:
//
//	tr, err := portwindow.New(30 * time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	tr.Observe(8080, time.Now())
//	w, ok := tr.Get(8080)
package portwindow
