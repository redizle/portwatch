// Package monitor ties together the scanner, state store, logger, and alerter
// into a single polling loop. Call New to create a Monitor from a config, then
// Run with a context to start watching. Cancel the context to stop cleanly.
//
// Typical usage:
//
//	cfg, _ := config.Load("portwatch.json")
//	log, _ := logger.New("stdout", "json")
//	alerter := alert.New(cfg.AlertHooks)
//	m := monitor.New(cfg, log, alerter)
//	if err := m.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatal(err)
//	}
package monitor

import "fmt"
