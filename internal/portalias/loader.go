package portalias

import (
	"encoding/json"
	"fmt"
	"os"
)

type aliasEntry struct {
	Port  int    `json:"port"`
	Label string `json:"label"`
}

// LoadFile reads a JSON file of alias entries into the Store.
// The file should contain an array of {"port": N, "label": "name"} objects.
func (s *Store) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portalias: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []aliasEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portalias: decode %s: %w", path, err)
	}

	for _, e := range entries {
		if err := s.Set(e.Port, e.Label); err != nil {
			return fmt.Errorf("portalias: entry port %d: %w", e.Port, err)
		}
	}
	return nil
}
