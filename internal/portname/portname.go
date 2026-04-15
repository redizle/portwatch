// Package portname resolves well-known port numbers to human-readable service names.
package portname

import "fmt"

// wellKnown maps common port numbers to their canonical service names.
var wellKnown = map[int]string{
	20:   "ftp-data",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Resolver resolves port numbers to service names, supporting custom overrides.
type Resolver struct {
	overrides map[int]string
}

// New returns a new Resolver with no overrides.
func New() *Resolver {
	return &Resolver{
		overrides: make(map[int]string),
	}
}

// SetOverride registers a custom name for the given port, taking precedence over
// the built-in well-known table.
func (r *Resolver) SetOverride(port int, name string) {
	r.overrides[port] = name
}

// Resolve returns the service name for the given port. Custom overrides are
// checked first, then the well-known table. If no match is found it returns
// a formatted fallback string like "port-8888".
func (r *Resolver) Resolve(port int) string {
	if name, ok := r.overrides[port]; ok {
		return name
	}
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return fmt.Sprintf("port-%d", port)
}

// IsWellKnown reports whether the port exists in the built-in table,
// regardless of any overrides.
func (r *Resolver) IsWellKnown(port int) bool {
	_, ok := wellKnown[port]
	return ok
}
