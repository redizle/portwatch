// Package digest provides periodic aggregation of port scan events into
// human-readable summaries. It is useful for batching alerts and producing
// interval-based reports without flooding downstream handlers.
//
// Usage:
//
//	b := digest.NewBuilder(5*time.Minute, resolver.Resolve)
//	b.Record(80, "open")
//	b.Record(443, "closed")
//	fmt.Println(b.Summary())
//	entries := b.Build() // flushes buffer
package digest
