package porttimelimit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonEntry struct {
	Port  int    `json:"port"`
	Start string `json:"start"` // e.g. "09:00"
	End   string `json:"end"`   // e.g. "17:00"
}

// LoadFile reads a JSON file of time-limit entries and populates the Limiter.
// Each entry must have a port, start time, and end time in "HH:MM" format.
func LoadFile(path string) (*Limiter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("porttimelimit: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []jsonEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("porttimelimit: decode %s: %w", path, err)
	}

	l := New()
	for _, e := range entries {
		start, err := parseClock(e.Start)
		if err != nil {
			return nil, fmt.Errorf("porttimelimit: port %d start: %w", e.Port, err)
		}
		end, err := parseClock(e.End)
		if err != nil {
			return nil, fmt.Errorf("porttimelimit: port %d end: %w", e.Port, err)
		}
		if err := l.Set(e.Port, Window{Start: start, End: end}); err != nil {
			return nil, err
		}
	}
	return l, nil
}

func parseClock(s string) (time.Duration, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return 0, fmt.Errorf("invalid time %q: %w", s, err)
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("time out of range: %q", s)
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}
