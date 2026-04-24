package portroute

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileEntry struct {
	Port  int    `json:"port"`
	Route string `json:"route"`
}

// LoadFile reads a JSON file of port-to-route mappings and populates r.
// The file should contain an array of {"port": N, "route": "..."} objects.
func LoadFile(path string, r *Router) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portroute: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portroute: decode %s: %w", path, err)
	}

	for _, e := range entries {
		if err := r.Set(e.Port, e.Route); err != nil {
			return fmt.Errorf("portroute: entry port %d: %w", e.Port, err)
		}
	}
	return nil
}
