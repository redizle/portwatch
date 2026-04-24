// Package portroute associates TCP port numbers with named route identifiers
// or service endpoint paths.
//
// Routes can be loaded from a JSON file or set programmatically. They are
// intended to provide human-readable context when displaying port activity
// in reports and alerts.
//
// Example JSON file format:
//
//	[
//	  {"port": 80,  "route": "/web"},
//	  {"port": 443, "route": "/secure"},
//	  {"port": 8080,"route": "/api/v1"}
//	]
package portroute
