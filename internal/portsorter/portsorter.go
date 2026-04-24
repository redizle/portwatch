// Package portsorter provides sorting utilities for port lists
// based on various criteria such as port number, status, priority, or label.
package portsorter

import (
	"sort"
)

// SortBy defines the field to sort by.
type SortBy int

const (
	ByPort SortBy = iota
	ByStatus
	ByLabel
)

// Order defines sort direction.
type Order int

const (
	Ascending Order = iota
	Descending
)

// Entry represents a sortable port entry.
type Entry struct {
	Port   int
	Status string
	Label  string
}

// Sorter sorts port entries by a configurable field and order.
type Sorter struct {
	field SortBy
	order Order
}

// New creates a new Sorter with the given field and order.
func New(field SortBy, order Order) *Sorter {
	return &Sorter{field: field, order: order}
}

// Sort sorts a slice of Entry values in place.
func (s *Sorter) Sort(entries []Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		var less bool
		switch s.field {
		case ByStatus:
			less = entries[i].Status < entries[j].Status
		case ByLabel:
			less = entries[i].Label < entries[j].Label
		default:
			less = entries[i].Port < entries[j].Port
		}
		if s.order == Descending {
			return !less
		}
		return less
	})
}

// SortedCopy returns a sorted copy of the given entries without modifying the original.
func (s *Sorter) SortedCopy(entries []Entry) []Entry {
	cp := make([]Entry, len(entries))
	copy(cp, entries)
	s.Sort(cp)
	return cp
}
