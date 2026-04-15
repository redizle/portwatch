package portname

import "testing"

func TestResolve_WellKnownHTTP(t *testing.T) {
	r := New()
	if got := r.Resolve(80); got != "http" {
		t.Errorf("expected http, got %s", got)
	}
}

func TestResolve_WellKnownSSH(t *testing.T) {
	r := New()
	if got := r.Resolve(22); got != "ssh" {
		t.Errorf("expected ssh, got %s", got)
	}
}

func TestResolve_Unknown_ReturnsFallback(t *testing.T) {
	r := New()
	if got := r.Resolve(9999); got != "port-9999" {
		t.Errorf("expected port-9999, got %s", got)
	}
}

func TestResolve_Override_TakesPrecedence(t *testing.T) {
	r := New()
	r.SetOverride(80, "my-web")
	if got := r.Resolve(80); got != "my-web" {
		t.Errorf("expected my-web, got %s", got)
	}
}

func TestResolve_Override_UnknownPort(t *testing.T) {
	r := New()
	r.SetOverride(12345, "custom-svc")
	if got := r.Resolve(12345); got != "custom-svc" {
		t.Errorf("expected custom-svc, got %s", got)
	}
}

func TestIsWellKnown_True(t *testing.T) {
	r := New()
	if !r.IsWellKnown(443) {
		t.Error("expected 443 to be well-known")
	}
}

func TestIsWellKnown_False(t *testing.T) {
	r := New()
	if r.IsWellKnown(9999) {
		t.Error("expected 9999 to not be well-known")
	}
}

func TestIsWellKnown_IgnoresOverride(t *testing.T) {
	// An override on an unknown port should not make it "well-known".
	r := New()
	r.SetOverride(9999, "custom")
	if r.IsWellKnown(9999) {
		t.Error("override should not affect IsWellKnown")
	}
}

func TestResolve_MultipleOverrides(t *testing.T) {
	r := New()
	r.SetOverride(3306, "primary-db")
	r.SetOverride(5432, "analytics-db")

	if got := r.Resolve(3306); got != "primary-db" {
		t.Errorf("expected primary-db, got %s", got)
	}
	if got := r.Resolve(5432); got != "analytics-db" {
		t.Errorf("expected analytics-db, got %s", got)
	}
}
