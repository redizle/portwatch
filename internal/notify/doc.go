// Package notify implements threshold-aware notification dispatch for portwatch.
//
// It supports multiple registered handlers and enforces a per-port cooldown
// to prevent notification floods when a port repeatedly changes state.
//
// Usage:
//
//	n := notify.New(30 * time.Second)
//	n.Register(func(e notify.Event) error {
//		fmt.Println(e.Message)
//		return nil
//	})
//	n.Dispatch(8080, "open", notify.LevelWarn)
package notify
