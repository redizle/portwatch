// Package portprofile provides a thread-safe store for associating named
// configuration profiles with individual ports.
//
// A Profile bundles common metadata — label, owner, priority, and tags —
// into a single value that can be attached to any valid port number.
// This is useful when multiple subsystems need to share the same
// descriptive context for a port without duplicating state.
//
// Example:
//
//	s := portprofile.New()
//	_ = s.Set(8080, portprofile.Profile{
//		Name:     "dev-api",
//		Label:    "Development API",
//		Owner:    "backend-team",
//		Priority: 2,
//		Tags:     []string{"internal", "dev"},
//	})
//	p, ok := s.Get(8080)
package portprofile
