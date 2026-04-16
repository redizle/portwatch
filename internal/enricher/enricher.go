// Package enricher attaches metadata to port events before they are
// dispatched or logged. It combines tagger and portname lookups into
// a single enrichment step.
package enricher

import (
	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/tagger"
)

// Meta holds enriched metadata for a single port.
type Meta struct {
	Port    int
	Name    string
	Tag     string
	Priority string
}

// Enricher attaches name and tag metadata to ports.
type Enricher struct {
	names  *portname.Resolver
	tags   *tagger.Tagger
}

// New returns a new Enricher using the given resolver and tagger.
// If either is nil a zero-value default is constructed.
func New(r *portname.Resolver, t *tagger.Tagger) *Enricher {
	if r == nil {
		r = portname.New(nil)
	}
	if t == nil {
		t = tagger.New(nil)
	}
	return &Enricher{names: r, tags: t}
}

// Enrich returns a Meta struct populated with the name and tag for port.
func (e *Enricher) Enrich(port int) Meta {
	return Meta{
		Port: port,
		Name: e.names.Resolve(port),
		Tag:  e.tags.Tag(port),
	}
}

// EnrichAll returns a slice of Meta for each port in ports.
func (e *Enricher) EnrichAll(ports []int) []Meta {
	out := make([]Meta, 0, len(ports))
	for _, p := range ports {
		out = append(out, e.Enrich(p))
	}
	return out
}
