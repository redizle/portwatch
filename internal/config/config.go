package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	Ports        []int         `json:"ports"`
	ScanInterval time.Duration `json:"scan_interval"`
	LogFile      string        `json:"log_file"`
	AlertHook    string        `json:"alert_hook"`
	AlertOnOpen  bool          `json:"alert_on_open"`
	AlertOnClose bool          `json:"alert_on_close"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() *Config {
	return &Config{
		Ports:        []int{80, 443, 8080, 8443},
		ScanInterval: 30 * time.Second,
		LogFile:      "portwatch.log",
		AlertOnOpen:  true,
		AlertOnClose: false,
	}
}

// Load reads a config from a JSON file at the given path.
// If the file doesn't exist, it returns the default config.
func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that the config values are sensible.
func (c *Config) Validate() error {
	if len(c.Ports) == 0 {
		return fmt.Errorf("at least one port must be specified")
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("port %d is out of valid range (1-65535)", p)
		}
	}
	if c.ScanInterval < time.Second {
		return fmt.Errorf("scan_interval must be at least 1 second")
	}
	return nil
}
