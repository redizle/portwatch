package monitor

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Monitor orchestrates scanning, state tracking, logging, and alerting.
type Monitor struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	store   *state.Store
	log     *logger.Logger
	alerter *alert.Alerter
}

// New creates a Monitor with all dependencies wired up.
func New(cfg *config.Config, log *logger.Logger, alerter *alert.Alerter) *Monitor {
	return &Monitor{
		cfg:     cfg,
		scanner: scanner.New(cfg.Timeout),
		store:   state.New(),
		log:     log,
		alerter: alerter,
	}
}

// Run starts the monitoring loop, blocking until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	m.log.Info("monitor started", map[string]interface{}{
		"ports":    m.cfg.Ports,
		"interval": m.cfg.Interval.String(),
	})

	for {
		select {
		case <-ctx.Done():
			m.log.Info("monitor stopped", nil)
			return ctx.Err()
		case <-ticker.C:
			m.scan()
		}
	}
}

func (m *Monitor) scan() {
	results := m.scanner.ScanPorts(m.cfg.Ports)
	for _, r := range results {
		changed := m.store.Update(r.Port, r.Open)
		if !changed {
			continue
		}
		fields := map[string]interface{}{
			"port": r.Port,
			"open": r.Open,
		}
		if r.Open {
			m.log.Info("port opened", fields)
		} else {
			m.log.Info("port closed", fields)
		}
		if err := m.alerter.Send(r.Port, r.Open); err != nil {
			m.log.Error("alert failed", map[string]interface{}{"error": err.Error()})
		}
	}
}
