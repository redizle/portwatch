package daemon

import (
	"fmt"

	"github.com/user/portwatch/internal/notify"
)

// HandlerFunc is a convenience type matching notify.HandlerFunc.
type HandlerFunc = notify.HandlerFunc

// Option is a functional option for configuring a Daemon after construction.
type Option func(*Daemon) error

// WithHandler registers an additional notification handler on the daemon's
// dispatcher. Handlers are called in registration order.
func WithHandler(name string, fn HandlerFunc) Option {
	return func(d *Daemon) error {
		if name == "" {
			return fmt.Errorf("handler name must not be empty")
		}
		d.notifier.Register(name, fn)
		return nil
	}
}

// Apply runs all options against the daemon, returning the first error
// encountered.
func (d *Daemon) Apply(opts ...Option) error {
	for _, o := range opts {
		if err := o(d); err != nil {
			return err
		}
	}
	return nil
}
