// Package daemon is the top-level orchestrator for portwatch.
//
// It wires the scanner, snapshot manager, filter, history, notifier, and
// reporter together into a single long-running process that is driven by a
// configurable polling interval.
//
// Typical usage:
//
//	cfg, _ := config.Load("portwatch.json")
//	log := slog.Default()
//	d, _ := daemon.New(cfg, log)
//	d.Run(ctx)
package daemon
