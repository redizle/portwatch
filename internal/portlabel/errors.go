package portlabel

import "errors"

// ErrInvalidPort is returned when a port number is out of range.
var ErrInvalidPort = errors.New("portlabel: port must be between 1 and 65535")
