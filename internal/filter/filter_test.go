package filter

import (
	"testing"
)

func TestNew_EmptyRules(t *testing.T) {
	f, err := New(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed with no rules")
	}
}

func TestNew_InvalidRule(t *testing.T) {
	_, err := New([]string{"notaport"}, nil)
	if err == nil {
		t.Error("expected error for invalid rule")
	}
}

func TestNew_OutOfRangeRule(t *testing.T) {
	_, err := New([]string{"0-100"}, nil)
	if err == nil {
		t.Error("expected error for out-of-range rule")
	}
}

func TestAllow_IncludeRange(t *testing.T) {
	f, err := New([]string{"8000-9000"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be blocked (not in include range)")
	}
}

func TestAllow_ExcludeRange(t *testing.T) {
	f, err := New(nil, []string{"1-1024"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow(443) {
		t.Error("expected 443 to be excluded")
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	f, err := New([]string{"80-9000"}, []string{"443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow(443) {
		t.Error("expected 443 to be excluded even though it's in include range")
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
}

func TestAllow_SinglePort(t *testing.T) {
	f, err := New([]string{"22"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(22) {
		t.Error("expected port 22 to be allowed")
	}
	if f.Allow(23) {
		t.Error("expected port 23 to be blocked")
	}
}
