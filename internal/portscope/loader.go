package portscope

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileEntry struct {
	Port  int    `json:"port"`
	Scope string `json:"scope"`
}

// LoadFile reads a JSON file of port scope assignments into c.
// Format: [{"port": 80, "scope": "external"}, ...]
func LoadFile(path string, c *Classifier) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portscope: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portscope: decode %s: %w", path, err)
	}

	for _, e := range entries {
		if err := c.Set(e.Port, Scope(e.Scope)); err != nil {
			return fmt.Errorf("portscope: entry port=%d: %w", e.Port, err)
		}
	}
	return nil
}
