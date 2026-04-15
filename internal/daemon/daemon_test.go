package daemon_test

import (
	"context"
	"testing"
	"time"

	"log/slog"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func testConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Ports = []int{19999} // unlikely to be open in CI
	cfg.Interval = 1
	cfg.MaxHistory = 10
	return cfg
}

func TestNew_ReturnsNonNil(t *testing.T) {
	cfg := testConfig()
	log := slog.Default()

	d, err := daemon.New(cfg, log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil daemon")
	}
}

func TestNew_InvalidFilterRule(t *testing.T) {
	cfg := testConfig()
	cfg.FilterRules = []string{"not-a-rule"}
	log := slog.Default()

	_, err := daemon.New(cfg, log)
	if err == nil {
		t.Fatal("expected error for invalid filter rule")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	cfg := testConfig()
	log := slog.Default()

	d, err := daemon.New(cfg, log)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got: %v", err)
	}
}

func TestRun_CompletesOneTick(t *testing.T) {
	cfg := testConfig()
	cfg.Interval = 60 // long interval — only the immediate first tick fires
	log := slog.Default()

	d, err := daemon.New(cfg, log)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	d.Run(ctx) //nolint:errcheck
	elapsed := time.Since(start)

	// should have blocked for ~200ms, not 60s
	if elapsed > 2*time.Second {
		t.Fatalf("run took too long: %v", elapsed)
	}
}
