package tagger

import (
	"testing"
)

func TestTag_WellKnown(t *testing.T) {
	tr := New(nil)
	cases := []struct {
		port int
		want string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{6379, "redis"},
	}
	for _, c := range cases {
		got := tr.Tag(c.port)
		if got != c.want {
			t.Errorf("Tag(%d) = %q, want %q", c.port, got, c.want)
		}
	}
}

func TestTag_Unknown(t *testing.T) {
	tr := New(nil)
	if got := tr.Tag(9999); got != "unknown" {
		t.Errorf("expected \"unknown\", got %q", got)
	}
}

func TestTag_OverrideTakesPrecedence(t *testing.T) {
	tr := New(map[int]string{80: "my-app"})
	if got := tr.Tag(80); got != "my-app" {
		t.Errorf("expected override \"my-app\", got %q", got)
	}
}

func TestSet_AddsOverride(t *testing.T) {
	tr := New(nil)
	tr.Set(9000, "custom-svc")
	if got := tr.Tag(9000); got != "custom-svc" {
		t.Errorf("expected \"custom-svc\", got %q", got)
	}
}

func TestRemove_FallsBackToWellKnown(t *testing.T) {
	tr := New(map[int]string{22: "my-ssh"})
	tr.Remove(22)
	if got := tr.Tag(22); got != "ssh" {
		t.Errorf("expected fallback \"ssh\", got %q", got)
	}
}

func TestRemove_FallsBackToUnknown(t *testing.T) {
	tr := New(map[int]string{9999: "temp"})
	tr.Remove(9999)
	if got := tr.Tag(9999); got != "unknown" {
		t.Errorf("expected \"unknown\", got %q", got)
	}
}

func TestNew_DoesNotMutateInput(t *testing.T) {
	input := map[int]string{80: "original"}
	tr := New(input)
	tr.Set(80, "changed")
	if input[80] != "original" {
		t.Error("New should copy the overrides map, not reference it")
	}
}
