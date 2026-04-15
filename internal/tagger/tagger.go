// Package tagger assigns human-readable tags to ports based on well-known
// service mappings and user-defined overrides.
package tagger

import "sync"

// wellKnown maps common port numbers to service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger resolves a tag string for a given port number.
type Tagger struct {
	mu        sync.RWMutex
	overrides map[int]string
}

// New returns a Tagger with an optional set of user overrides.
func New(overrides map[int]string) *Tagger {
	ov := make(map[int]string, len(overrides))
	for k, v := range overrides {
		ov[k] = v
	}
	return &Tagger{overrides: ov}
}

// Tag returns the tag for the given port. User overrides take precedence
// over well-known mappings. Returns "unknown" when no match is found.
func (t *Tagger) Tag(port int) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if tag, ok := t.overrides[port]; ok {
		return tag
	}
	if tag, ok := wellKnown[port]; ok {
		return tag
	}
	return "unknown"
}

// Set adds or updates a user override for the given port.
func (t *Tagger) Set(port int, tag string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.overrides[port] = tag
}

// Remove deletes a user override, falling back to well-known or "unknown".
func (t *Tagger) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.overrides, port)
}
