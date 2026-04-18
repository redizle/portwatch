// Package portflag provides a simple mechanism for flagging ports with
// a short reason string, useful for marking ports that require manual
// review or investigation during monitoring sessions.
//
// Flags are stored in memory and are not persisted across restarts.
// Use Set to flag a port, Unflag to clear it, and All to enumerate
// all currently flagged ports.
package portflag
