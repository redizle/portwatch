// Package porttag provides a simple store for attaching and querying
// string tags on port numbers.
//
// Tags are arbitrary labels (e.g. "web", "internal", "monitored") that
// can be associated with any valid port (1–65535). Multiple tags per
// port are supported. All operations are safe for concurrent use.
package porttag
