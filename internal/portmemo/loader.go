package portmemo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type fileEntry struct {
	Port int    `json:"port"`
	Text string `json:"text"`
	TTL  int    `json:"ttl_seconds"`
}

// LoadFile reads a JSON file of memo entries and populates the store.
// Each entry may include an optional ttl_seconds field (0 = no expiry).
func LoadFile(path string, s *Store) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("portmemo: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []fileEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("portmemo: decode %s: %w", path, err)
	}

	for _, e := range entries {
		ttl := time.Duration(e.TTL) * time.Second
		if err := s.Set(e.Port, e.Text, ttl); err != nil {
			return fmt.Errorf("portmemo: entry port %d: %w", e.Port, err)
		}
	}
	return nil
}
