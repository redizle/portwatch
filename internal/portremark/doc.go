// Package portremark provides a thread-safe store for attaching
// short human-readable remarks to monitored ports.
//
// Remarks are timestamped at creation and can be retrieved,
// overwritten, or removed at any time. They are intended for
// operator notes such as "pending firewall change" or "under review".
package portremark
