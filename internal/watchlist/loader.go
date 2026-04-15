package watchlist

import (
	"encoding/json"
	"fmt"
	"os"
)

// jsonEntry mirrors Entry for JSON unmarshalling.
type jsonEntry struct {
	Port     int    `json:"port"`
	Label    string `json:"label"`
	Priority int    `json:"priority"`
}

// LoadFile reads a JSON file containing an array of port entries and
// populates the watchlist. Existing entries are preserved.
//
// Example file format:
//
//	[
//	  {"port": 22,   "label": "ssh",   "priority": 3},
//	  {"port": 8080, "label": "dev",   "priority": 1}
//	]
func (w *Watchlist) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("watchlist: open %s: %w", path, err)
	}
	defer f.Close()

	var raw []jsonEntry
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return fmt.Errorf("watchlist: decode %s: %w", path, err)
	}

	for _, r := range raw {
		pri := Priority(r.Priority)
		if pri < PriorityLow || pri > PriorityHigh {
			pri = PriorityNormal
		}
		if err := w.Add(r.Port, r.Label, pri); err != nil {
			return err
		}
	}
	return nil
}
