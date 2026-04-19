package portdesc

import (
	"testing"
)

func TestResolve_WellKnownHTTP(t *testing.T) {
	r := New()
	if got := r.Resolve(80); got != "HTTP" {
		t.Fatalf("expected HTTP, got %s", got)
	}
}

func TestResolve_WellKnownSSH(t *testing.T) {
	r := New()
	if got := r.Resolve(22); got != "SSH" {
		t.Fatalf("expected SSH, got %s", got)
	}
}

func TestResolve_Unknown(t *testing.T) {
	r := New()
	if got := r.Resolve(19999); got != "unknown" {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestSet_OverrideTakesPrecedence(t *testing.T) {
	r := New()
	if err := r.Set(80, "my-web"); err != nil {
		t.Fatal(err)
	}
	if got := r.Resolve(80); got != "my-web" {
		t.Fatalf("expected my-web, got %s", got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Set(0, "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestSet_EmptyDesc(t *testing.T) {
	r := New()
	if err := r.Set(80, ""); err == nil {
		t.Fatal("expected error for empty description")
	}
}

func TestRemove_FallsBackToWellKnown(t *testing.T) {
	r := New()
	_ = r.Set(22, "custom-ssh")
	r.Remove(22)
	if got := r.Resolve(22); got != "SSH" {
		t.Fatalf("expected SSH after remove, got %s", got)
	}
}
