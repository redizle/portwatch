// Package portannot stores arbitrary key-value annotations against
// monitored ports. Annotations are intended for human-readable metadata
// such as environment labels, team ownership, or deployment tier.
//
// Usage:
//
//	s := portannot.New()
//	_ = s.Set(8080, "env", "staging")
//	v, ok := s.Get(8080, "env")
//
// Annotations can also be pre-loaded from a JSON file via LoadFile.
package portannot
