// Package daemon wires together all portwatch subsystems and runs the
// main scan/notify loop as a long-lived process.
package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon coordinates all subsystems.
type Daemon struct {
	cfg      *config.Config
	mon      *monitor.Monitor
	notifier *notify.Dispatcher
	history  *history.History
	reporter *report.Reporter
	manager  *snapshot.Manager
	filter   *filter.Filter
	log      *slog.Logger
}

// New constructs a Daemon from the provided config and logger.
func New(cfg *config.Config, log *slog.Logger) (*Daemon, error) {
	f, err := filter.New(cfg.FilterRules)
	if err != nil {
		return nil, err
	}

	h := history.New(cfg.MaxHistory, cfg.HistoryFile)
	r := report.New(h, log)
	n := notify.New(cfg.CooldownSeconds, log)
	sm := snapshot.NewManager()
	m := monitor.New(cfg, log)

	return &Daemon{
		cfg:      cfg,
		mon:      m,
		notifier: n,
		history:  h,
		reporter: r,
		manager:  sm,
		filter:   f,
		log:      log,
	}, nil
}

// Run starts the daemon and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	d.log.Info("portwatch daemon starting",
		"ports", len(d.cfg.Ports),
		"interval", d.cfg.Interval,
	)

	ticker := time.NewTicker(time.Duration(d.cfg.Interval) * time.Second)
	defer ticker.Stop()

	if err := d.tick(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			d.log.Info("daemon shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(ctx); err != nil {
				d.log.Error("tick error", "err", err)
			}
		}
	}
}

func (d *Daemon) tick(ctx context.Context) error {
	snap, err := d.mon.Scan(ctx)
	if err != nil {
		return err
	}

	d.manager.Rotate(snap)
	diffs := d.manager.Diff()

	for _, ev := range diffs {
		if !d.filter.Allow(ev.Port) {
			continue
		}
		d.history.Record(ev)
		d.notifier.Dispatch(ctx, ev)
	}
	return nil
}
