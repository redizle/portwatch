// Package portrelay provides a fan-out relay for port events.
//
// Targets are registered by name and receive every dispatched
// (port, event) pair. Errors from individual targets are collected
// and returned together so that a single failing target does not
// prevent delivery to the remaining ones.
//
// Example:
//
//	r := portrelay.New()
//	_ = r.Register(portrelay.Target{
//		Name: "logger",
//		Handler: func(port int, event string) error {
//			fmt.Printf("port %d: %s\n", port, event)
//			return nil
//		},
//	})
//	_ = r.Dispatch(8080, "open")
package portrelay
