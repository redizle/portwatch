// Package snapshot provides point-in-time capture and diffing of port states
// observed during a scan cycle.
//
// A Snapshot records which ports were seen as open or closed at a given moment.
// The Manager maintains a rolling pair of snapshots (current and previous),
// enabling efficient diff computation at the end of each scan cycle to surface
// only ports whose status has changed.
//
// Typical usage:
//
//	mgr := snapshot.NewManager()
//
//	// during a scan cycle:
//	for _, port := range scannedPorts {
//		mgr.Current().Set(port, isOpen)
//	}
//
//	// at end of cycle:
//	changed := mgr.Rotate()
//	for _, ps := range changed {
//		// handle open/close transition
//	}
package snapshot
