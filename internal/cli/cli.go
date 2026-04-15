// Package cli provides the command-line interface for portwatch.
package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

// Run parses CLI flags, loads config, and starts the daemon.
// It blocks until an interrupt signal is received or the context is cancelled.
func Run(args []string) int {
	fs := flag.NewFlagSet("portwatch", flag.ContinueOnError)

	configPath := fs.String("config", "", "path to config file (JSON)")
	verbose := fs.Bool("verbose", false, "enable verbose logging")
	version := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		return 1
	}

	if *version {
		fmt.Println("portwatch v0.1.0")
		return 0
	}

	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
			return 1
		}
	} else {
		cfg = config.DefaultConfig()
	}

	if *verbose {
		cfg.LogLevel = "debug"
	}

	d, err := daemon.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialise daemon: %v\n", err)
		return 1
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := d.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
		return 1
	}

	return 0
}
