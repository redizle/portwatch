package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/cli"
)

func TestRun_Version(t *testing.T) {
	code := cli.Run([]string{"-version"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRun_InvalidFlag(t *testing.T) {
	code := cli.Run([]string{"-unknownflag"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for unknown flag, got %d", code)
	}
}

func TestRun_MissingConfigFile(t *testing.T) {
	code := cli.Run([]string{"-config", "/nonexistent/path/config.json"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for missing config, got %d", code)
	}
}

func TestRun_ValidConfigFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	content := `{
		"ports": [19999],
		"interval_ms": 500,
		"log_format": "text"
	}`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Run with a very short-lived context — we just want no crash on startup.
	// The daemon will be cancelled immediately via the default context timeout
	// in tests; here we verify the exit code reflects a clean load.
	doneCh := make(chan int, 1)
	go func() {
		doneCh <- cli.Run([]string{"-config", cfgPath})
	}()

	// Send interrupt to self so the daemon exits cleanly.
	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(os.Interrupt)

	code := <-doneCh
	if code != 0 {
		t.Fatalf("expected clean exit code 0, got %d", code)
	}
}

func TestRun_DefaultConfig_NoArgs(t *testing.T) {
	doneCh := make(chan int, 1)
	go func() {
		doneCh <- cli.Run([]string{})
	}()

	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(os.Interrupt)

	code := <-doneCh
	if code != 0 {
		t.Fatalf("expected exit code 0 with default config, got %d", code)
	}
}
