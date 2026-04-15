package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Config holds the full portwatch daemon configuration.
type Config struct {
	Ports        []int         `json:"ports"`
	Interval     time.Duration `json:"interval"`
	LogFormat    string        `json:"log_format"`
	LogOutput    string        `json:"log_output"`
	AlertHooks   []string      `json:"alert_hooks"`
	HistoryFile  string        `json:"history_file"`
	MaxHistory   int           `json:"max_history"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Ports:      []int{80, 443, 8080},
		Interval:   10 * time.Second,
		LogFormat:  "text",
		LogOutput:  "stdout",
		AlertHooks: []string{},
		HistoryFile: "",
		MaxHistory: 1000,
	}
}

// Load reads a JSON config file from path, falling back to defaults for
// unspecified fields.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are acceptable.
func (c *Config) Validate() error {
	for := range c.Ports {
		if p < 1 || p > 65535 {
			return errors.New("config: port out of range (1-65535)")
		}
	}
	if c.Interval < time.Second {
		return errors.New("config: interval must be at least 1s")
	}
	if c.LogFormat != "text" && c.LogFormat != "json" {
		return errors.New("config: log_format must be 'text' or 'json'")
	}
	if c.MaxHistory < 0 {
		return errors.New("config: max_history must be non-negative")
	}
	return nil
}
