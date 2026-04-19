package portnote

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileEntry struct {
	Port int    `json:"port"`
	Text string `json:"text"`
}

// LoadFile reads a JSON file of port notes and populates s.
// The file must contain an array of {"port": N, "text": "..."} objects.
func LoadFile(s *Store, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portnote: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portnote: decode %s: %w", path, err)
	}

	for _, e := range entries {
		if err := s.Set(e.Port, e.Text); err != nil {
			return fmt.Errorf("portnote: entry port %d: %w", e.Port, err)
		}
	}
	return nil
}
