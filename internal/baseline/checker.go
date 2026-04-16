package baseline

import "fmt"

// Violation describes a port that is open but not baselined.
type Violation struct {
	Port    int
	Message string
}

// String implements fmt.Stringer.
func (v Violation) String() string {
	return fmt.Sprintf("port %d is open but not in baseline", v.Port)
}

// Checker wraps a Baseline and produces Violations.
type Checker struct {
	b *Baseline
}

// NewChecker returns a Checker backed by the given Baseline.
func NewChecker(b *Baseline) *Checker {
	return &Checker{b: b}
}

// Check compares active ports against the baseline and returns any violations.
func (c *Checker) Check(active []int) []Violation {
	unexpected := c.b.Unexpected(active)
	if len(unexpected) == 0 {
		return nil
	}
	out := make([]Violation, 0, len(unexpected))
	for _, p := range unexpected {
		out = append(out, Violation{
			Port:    p,
			Message: fmt.Sprintf("port %d is open but not in baseline", p),
		})
	}
	return out
}

// HasViolations returns true if any active ports are outside the baseline.
func (c *Checker) HasViolations(active []int) bool {
	return len(c.b.Unexpected(active)) > 0
}
