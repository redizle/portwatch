// Package portrank provides a scoring mechanism that ranks ports by activity
// frequency and alert weight.
//
// Scores are accumulated via Add and can be pinned to a fixed value with
// SetOverride. The Reporter prints a descending-score table suitable for
// operator dashboards.
//
// Example usage:
//
//	r := portrank.New()
//	_ = r.Add(80, 1)           // increment on each scan hit
//	_ = r.SetOverride(22, 100) // always rank SSH at the top
//	e, _ := r.Get(80)
//	fmt.Println(e.Score)
package portrank
