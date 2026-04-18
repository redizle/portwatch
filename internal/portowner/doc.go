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
//	owner, ok := r.Get(8080)
package portowner
