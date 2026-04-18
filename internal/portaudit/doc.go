// Package portaudit provides an in-memory, capped audit log for port
// lifecycle events within portwatch.
//
// Each call to Record appends a timestamped Entry describing what
// happened (opened, closed, suppressed, alerted) and an optional
// free-text note. Entries can be queried by port or exported as JSON
// or plain text via Exporter.
//
// Usage:
//
//	l := portaudit.New(500) // keep last 500 entries
//	l.Record(8080, portaudit.ActionOpened, "first seen")
//	for _, e := range l.ForPort(8080) {
//		fmt.Println(e)
//	}
package portaudit
