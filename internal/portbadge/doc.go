// Package portbadge provides a simple store for assigning icon and label badges
// to monitored ports. Badges are short human-readable status indicators used
// when rendering port summaries in the CLI or reports.
//
// Example:
//
//	s := portbadge.New()
//	s.Set(80, "✔", "http")
//	b := s.Get(80) // Badge{Icon:"✔", Label:"http"}
package portbadge
