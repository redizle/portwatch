// Package portowner provides a thread-safe registry for mapping port numbers
// to logical owner names or team labels.
//
// Owners can be used in alerts and reports to attribute port activity to
// specific teams or services.
//
// Example:
//
//	r := portowner.New()
//	r.Set(8080, "team-backend")
//	r.Set(9090, "team-metrics")
//	owner, ok := r.Get(8080)
//	// owner == "team-backend", ok == true
//
//	_, ok = r.Get(1234)
//	// ok == false, port has no registered owner
package portowner
