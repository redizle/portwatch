// Package portchain implements a middleware-style handler chain for port events.
//
// Handlers are executed in registration order. The chain halts on the first
// error returned by a handler, allowing early termination for filtering,
// suppression, or alerting logic.
//
// Example:
//
//	c, _ := portchain.NewBuilder().
//		Add(loggingHandler).
//		Add(alertHandler).
//		Build()
//
//	if err := c.Run(8080, "open"); err != nil {
//		log.Println("chain halted:", err)
//	}
package portchain
