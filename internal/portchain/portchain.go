// Package portchain provides a middleware-style processing pipeline for port events.
package portchain

import "errors"

// Handler processes a port event and returns an error to halt the chain.
type Handler func(port int, status string) error

// Chain executes a sequence of handlers for a port event.
type Chain struct {
	handlers []Handler
}

// New returns an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Use appends a handler to the chain.
func (c *Chain) Use(h Handler) error {
	if h == nil {
		return errors.New("portchain: handler must not be nil")
	}
	c.handlers = append(c.handlers, h)
	return nil
}

// Run executes all handlers in order for the given port and status.
// It stops and returns the first error encountered.
func (c *Chain) Run(port int, status string) error {
	if port < 1 || port > 65535 {
		return errors.New("portchain: port out of range")
	}
	for _, h := range c.handlers {
		if err := h(port, status); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of registered handlers.
func (c *Chain) Len() int {
	return len(c.handlers)
}

// Reset removes all handlers from the chain.
func (c *Chain) Reset() {
	c.handlers = nil
}
