// Package filter implements include/exclude filtering for port ranges
// in portwatch. It allows fine-grained control over which ports are
// actively monitored by the daemon.
//
// Rules are expressed as single ports ("80") or hyphen-separated ranges
// ("8000-9000"). Exclude rules always take precedence over include rules.
// If no include rules are configured, all ports are monitored unless
// explicitly excluded.
//
// Example usage:
//
//	f, err := filter.New(
//		[]string{"1000-9000"},  // only watch these
//		[]string{"3306", "5432"}, // but never these
//	)
//	if f.Allow(8080) { /* monitor it */ }
package filter
