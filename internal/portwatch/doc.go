// Package portwatch provides a unified per-port watch record that combines
// liveness tracking, labeling, and ownership into a single in-memory store.
//
// Use Touch to record each scan observation, SetLabel and SetOwner to attach
// metadata, and Get or All to query the current state.
package portwatch
