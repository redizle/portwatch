// Package filter provides port filtering logic for portwatch.
// It allows users to include or exclude specific ports or ranges
// from monitoring based on configured rules.
package filter

import "fmt"

// Rule defines a single filter rule.
type Rule struct {
	Low  int
	High int
	Mode string // "include" or "exclude"
}

// Filter holds the compiled filter rules.
type Filter struct {
	rules []Rule
}

// New creates a Filter from include and exclude port range strings.
// Ranges are expressed as "low-high" or single port numbers.
func New(includes, excludes []string) (*Filter, error) {
	f := &Filter{}
	for _, s := range includes {
		r, err := parseRule(s, "include")
		if err != nil {
			return nil, err
		}
		f.rules = append(f.rules, r)
	}
	for _, s := range excludes {
		r, err := parseRule(s, "exclude")
		if err != nil {
			return nil, err
		}
		f.rules = append(f.rules, r)
	}
	return f, nil
}

// Allow returns true if the given port should be monitored.
// Exclude rules take precedence over include rules.
// If no include rules exist, all ports are allowed unless excluded.
func (f *Filter) Allow(port int) bool {
	hasIncludes := false
	included := false
	for _, r := range f.rules {
		if r.Mode == "include" {
			hasIncludes = true
			if port >= r.Low && port <= r.High {
				included = true
			}
		}
		if r.Mode == "exclude" && port >= r.Low && port <= r.High {
			return false
		}
	}
	if hasIncludes {
		return included
	}
	return true
}

func parseRule(s, mode string) (Rule, error) {
	var low, high int
	_, err := fmt.Sscanf(s, "%d-%d", &low, &high)
	if err != nil {
		_, err = fmt.Sscanf(s, "%d", &low)
		if err != nil {
			return Rule{}, fmt.Errorf("invalid port rule %q", s)
		}
		high = low
	}
	if low < 1 || high > 65535 || low > high {
		return Rule{}, fmt.Errorf("port rule %q out of valid range", s)
	}
	return Rule{Low: low, High: high, Mode: mode}, nil
}
