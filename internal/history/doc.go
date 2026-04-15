// Package history provides an append-only record of port status change events
// observed by the portwatch monitor.
//
// Events are stored in memory and can optionally be persisted to a JSON file
// so that history survives daemon restarts. The maximum number of retained
// entries is configurable; older entries are evicted when the cap is reached.
//
// Typical usage:
//
//	h := history.New("/var/lib/portwatch/history.json", 500)
//	h.Record(8080, "open")
//	entries := h.GetByPort(8080)
package history
