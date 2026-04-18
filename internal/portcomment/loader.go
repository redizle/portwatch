package portcomment

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileEntry struct {
	Port    int    `json:"port"`
	Comment string `json:"comment"`
}

// LoadFile reads a JSON file of port comments into the Store.
// The file should contain an array of {"port": N, "comment": "..."} objects.
func (s *Store) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portcomment: open file: %w", err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portcomment: decode: %w", err)
	}

	for _, e := range entries {
		if err := s.Set(e.Port, e.Comment); err != nil {
			return fmt.Errorf("portcomment: invalid entry port %d: %w", e.Port, err)
		}
	}
	return nil
}
