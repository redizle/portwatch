// Package portmatch provides pattern-based port matching using glob-style rules.
package portmatch

import (
	"fmt"
	"strconv"
	"strings"
)

// Matcher holds compiled port match patterns.
type Matcher struct {
	rules []rule
}

type rule struct {
	raw   string
	low   int
	high  int
	exact int
	range_ bool
}

// New creates a Matcher from a slice of pattern strings.
// Patterns may be exact ports ("80"), ranges ("8000-8999"), or wildcard ("*").
func New(patterns []string) (*Matcher, error) {
	m := &Matcher{}
	for _, p := range patterns {
		r, err := parsePattern(p)
		if err != nil {
			return nil, fmt.Errorf("portmatch: invalid pattern %q: %w", p, err)
		}
		m.rules = append(m.rules, r)
	}
	return m, nil
}

// Match reports whether port matches any of the configured patterns.
func (m *Matcher) Match(port int) bool {
	for _, r := range m.rules {
		if r.raw == "*" {
			return true
		}
		if r.range_ {
			if port >= r.low && port <= r.high {
				return true
			}
		} else {
			if port == r.exact {
				return true
			}
		}
	}
	return false
}

// Len returns the number of patterns loaded.
func (m *Matcher) Len() int { return len(m.rules) }

func parsePattern(p string) (rule, error) {
	p = strings.TrimSpace(p)
	if p == "*" {
		return rule{raw: "*"}, nil
	}
	if strings.Contains(p, "-") {
		parts := strings.SplitN(p, "-", 2)
		lo, err1 := strconv.Atoi(parts[0])
		hi, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return rule{}, fmt.Errorf("bad range")
		}
		if lo < 1 || hi > 65535 || lo > hi {
			return rule{}, fmt.Errorf("range out of bounds")
		}
		return rule{raw: p, range_: true, low: lo, high: hi}, nil
	}
	n, err := strconv.Atoi(p)
	if err != nil || n < 1 || n > 65535 {
		return rule{}, fmt.Errorf("bad port number")
	}
	return rule{raw: p, exact: n}, nil
}
