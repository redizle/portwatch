package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Ports) == 0 {
		t.Error("expected default ports to be non-empty")
	}
	if cfg.ScanInterval < time.Second {
		t.Errorf("expected scan interval >= 1s, got %v", cfg.ScanInterval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	cfg, err := Load("/tmp/portwatch_nonexistent_config.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	input := &Config{
		Ports:        []int{9000, 9001},
		ScanInterval: 10 * time.Second,
		LogFile:      "test.log",
		AlertOnOpen:  true,
		AlertOnClose: true,
	}

	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(input); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Ports) != 2 || cfg.Ports[0] != 9000 {
		t.Errorf("unexpected ports: %v", cfg.Ports)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Ports = []int{0, 99999}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for out-of-range ports")
	}
}

func TestValidate_ShortInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ScanInterval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for interval < 1s")
	}
}

func TestValidate_NoPorts(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Ports = []int{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty ports list")
	}
}
