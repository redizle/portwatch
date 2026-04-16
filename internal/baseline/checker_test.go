package baseline

import (
	"strings"
	"testing"
)

func newChecker(ports ...int) *Checker {
	b := New()
	for _, p := range ports {
		b.Add(p, "")
	}
	return NewChecker(b)
}

func TestCheck_NoViolations(t *testing.T) {
	c := newChecker(80, 443)
	v := c.Check([]int{80, 443})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_WithViolations(t *testing.T) {
	c := newChecker(80)
	v := c.Check([]int{80, 8080})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %d", v[0].Port)
	}
}

func TestCheck_EmptyActive(t *testing.T) {
	c := newChecker(80)
	v := c.Check([]int{})
	if len(v) != 0 {
		t.Fatal("expected no violations for empty active list")
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Port: 3306, Message: "port 3306 is open but not in baseline"}
	if !strings.Contains(v.String(), "3306") {
		t.Fatal("expected port number in violation string")
	}
}

func TestHasViolations_True(t *testing.T) {
	c := newChecker(80)
	if !c.HasViolations([]int{80, 9999}) {
		t.Fatal("expected violations")
	}
}

func TestHasViolations_False(t *testing.T) {
	c := newChecker(80, 443)
	if c.HasViolations([]int{80, 443}) {
		t.Fatal("expected no violations")
	}
}
