// Package portmetadata provides a thread-safe key-value metadata store
// for attaching arbitrary string fields to port numbers.
//
// Each port can hold multiple named metadata entries. Entries can be
// individually deleted or cleared in bulk. The store is safe for
// concurrent use.
//
// Example:
//
//	s := portmetadata.New()
//	_ = s.Set(8080, "env", "production")
//	v, ok := s.Get(8080, "env")
package portmetadata
