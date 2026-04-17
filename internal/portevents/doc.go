// Package portevents implements a lightweight publish/subscribe event bus
// for port state changes within portwatch.
//
// Consumers subscribe to specific EventType values (opened, closed, changed)
// and receive Event structs whenever the daemon detects a matching transition.
//
// The Bus is safe for concurrent use.
package portevents
