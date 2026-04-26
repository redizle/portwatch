// Package portclassify assigns a classification tier to ports based on
// their number range: system (1–1023), registered (1024–49151), or
// dynamic (49152–65535).
package portclassify

import (
	"fmt"
	"sync"
)

// Tier represents a port classification tier.
type Tier string

const (
	TierSystem     Tier = "system"
	TierRegistered Tier = "registered"
	TierDynamic    Tier = "dynamic"
	TierUnknown    Tier = "unknown"
)

// String implements fmt.Stringer.
func (t Tier) String() string { return string(t) }

// Classifier maps ports to classification tiers, with optional overrides.
type Classifier struct {
	mu        sync.RWMutex
	overrides map[int]Tier
}

// New returns a new Classifier with no overrides.
func New() *Classifier {
	return &Classifier{overrides: make(map[int]Tier)}
}

// Classify returns the Tier for the given port number.
// Override entries take precedence over the default range-based logic.
func (c *Classifier) Classify(port int) (Tier, error) {
	if port < 1 || port > 65535 {
		return TierUnknown, fmt.Errorf("portclassify: port %d out of range", port)
	}
	c.mu.RLock()
	if t, ok := c.overrides[port]; ok {
		c.mu.RUnlock()
		return t, nil
	}
	c.mu.RUnlock()
	return defaultTier(port), nil
}

// Override sets a custom Tier for a specific port.
func (c *Classifier) Override(port int, tier Tier) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portclassify: port %d out of range", port)
	}
	c.mu.Lock()
	c.overrides[port] = tier
	c.mu.Unlock()
	return nil
}

// ClearOverride removes any custom override for the given port.
func (c *Classifier) ClearOverride(port int) {
	c.mu.Lock()
	delete(c.overrides, port)
	c.mu.Unlock()
}

func defaultTier(port int) Tier {
	switch {
	case port <= 1023:
		return TierSystem
	case port <= 49151:
		return TierRegistered
	default:
		return TierDynamic
	}
}
