package portmatch_test

import (
	"testing"

	"portwatch/internal/portmatch"
)

func TestNew_EmptyPatterns(t *testing.T) {
	m, err := portmatch.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 0 {
		t.Errorf("expected 0 rules, got %d", m.Len())
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := portmatch.New([]string{"abc"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_OutOfRangePort(t *testing.T) {
	_, err := portmatch.New([]string{"99999"})
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestMatch_ExactPort(t *testing.T) {
	m, _ := portmatch.New([]string{"80"})
	if !m.Match(80) {
		t.Error("expected port 80 to match")
	}
	if m.Match(81) {
		t.Error("expected port 81 not to match")
	}
}

func TestMatch_Range(t *testing.T) {
	m, _ := portmatch.New([]string{"8000-8080"})
	if !m.Match(8000) {
		t.Error("expected 8000 to match")
	}
	if !m.Match(8080) {
		t.Error("expected 8080 to match")
	}
	if m.Match(7999) {
		t.Error("expected 7999 not to match")
	}
	if m.Match(8081) {
		t.Error("expected 8081 not to match")
	}
}

func TestMatch_Wildcard(t *testing.T) {
	m, _ := portmatch.New([]string{"*"})
	for _, p := range []int{1, 80, 443, 65535} {
		if !m.Match(p) {
			t.Errorf("expected port %d to match wildcard", p)
		}
	}
}

func TestMatch_MultiplePatterns(t *testing.T) {
	m, _ := portmatch.New([]string{"22", "443", "9000-9100"})
	for _, p := range []int{22, 443, 9050} {
		if !m.Match(p) {
			t.Errorf("expected port %d to match", p)
		}
	}
	if m.Match(80) {
		t.Error("expected 80 not to match")
	}
}

func TestNew_InvalidRange(t *testing.T) {
	_, err := portmatch.New([]string{"9000-8000"})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}
