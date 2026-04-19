// Package portdesc provides short human-readable descriptions for ports.
package portdesc

import "fmt"

// Descriptions maps well-known ports to a short description.
var wellKnown = map[int]string{
	21:   "FTP",
	22:   "SSH",
	23:   "Telnet",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	110:  "POP3",
	143:  "IMAP",
	443:  "HTTPS",
	3306: "MySQL",
	5432: "PostgreSQL",
	6379: "Redis",
	8080: "HTTP-alt",
	27017: "MongoDB",
}

// Resolver resolves descriptions for ports.
type Resolver struct {
	overrides map[int]string
}

// New returns a new Resolver.
func New() *Resolver {
	return &Resolver{overrides: make(map[int]string)}
}

// Set adds or replaces a description for a port.
func (r *Resolver) Set(port int, desc string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portdesc: invalid port %d", port)
	}
	if desc == "" {
		return fmt.Errorf("portdesc: description must not be empty")
	}
	r.overrides[port] = desc
	return nil
}

// Resolve returns the description for a port. Overrides take precedence over
// well-known defaults. Falls back to "unknown" if nothing is found.
func (r *Resolver) Resolve(port int) string {
	if d, ok := r.overrides[port]; ok {
		return d
	}
	if d, ok := wellKnown[port]; ok {
		return d
	}
	return "unknown"
}

// Remove deletes an override for a port.
func (r *Resolver) Remove(port int) {
	delete(r.overrides, port)
}
