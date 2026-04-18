package portannot

import (
	"encoding/json"
	"fmt"
	"os"
)

// fileEntry is the JSON shape for a single port's annotations.
type fileEntry struct {
	Port        int               `json:"port"`
	Annotations map[string]string `json:"annotations"`
}

// LoadFile reads a JSON file and populates the store.
// The file should contain an array of {port, annotations} objects.
func LoadFile(path string, s *Store) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portannot: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portannot: decode %s: %w", path, err)
	}

	for _, e := range entries {
		for k, v := range e.Annotations {
			if err := s.Set(e.Port, k, v); err != nil {
				return fmt.Errorf("portannot: port %d key %q: %w", e.Port, k, err)
			}
		}
	}
	return nil
}
