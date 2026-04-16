package digest

import "time"

// Builder accumulates port change events and produces a Digest.
type Builder struct {
	digest *Digest
	namer  func(int) string
}

// NewBuilder returns a Builder using the provided window and name resolver.
func NewBuilder(window time.Duration, namer func(int) string) *Builder {
	if namer == nil {
		namer = func(p int) string { return "unknown" }
	}
	return &Builder{
		digest: New(window),
		namer:  namer,
	}
}

// Record adds a port status change to the builder.
func (b *Builder) Record(port int, status string) {
	b.digest.Add(Entry{
		Port:   port,
		Status: status,
		Name:   b.namer(port),
	})
}

// Build flushes and returns the current digest entries.
func (b *Builder) Build() []Entry {
	return b.digest.Flush()
}

// Summary returns the digest summary without flushing.
func (b *Builder) Summary() string {
	return b.digest.Summary()
}

// Len returns the number of pending entries.
func (b *Builder) Len() int {
	return b.digest.Len()
}
