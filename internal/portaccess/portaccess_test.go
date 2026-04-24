package portaccess_test

import (
	"testing"

	"portwatch/internal/portaccess"
)

func TestSet_And_Get(t *testing.T) {
	tr := portaccess.New()
	if err := tr.Set(8080, portaccess.PolicyAllow, "trusted service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Policy != portaccess.PolicyAllow {
		t.Errorf("expected allow, got %v", e.Policy)
	}
	if e.Reason != "trusted service" {
		t.Errorf("unexpected reason: %s", e.Reason)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	tr := portaccess.New()
	if err := tr.Set(0, portaccess.PolicyDeny, "bad"); err == nil {
		t.Error("expected error for port 0")
	}
	if err := tr.Set(70000, portaccess.PolicyDeny, "bad"); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestGet_Missing(t *testing.T) {
	tr := portaccess.New()
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected no entry for unknown port")
	}
}

func TestRemove(t *testing.T) {
	tr := portaccess.New()
	_ = tr.Set(443, portaccess.PolicyDeny, "blocked")
	tr.Remove(443)
	_, ok := tr.Get(443)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestIsAllowed_DefaultsToTrue(t *testing.T) {
	tr := portaccess.New()
	if !tr.IsAllowed(1234) {
		t.Error("expected unknown port to default to allowed")
	}
}

func TestIsAllowed_Deny(t *testing.T) {
	tr := portaccess.New()
	_ = tr.Set(22, portaccess.PolicyDeny, "ssh blocked")
	if tr.IsAllowed(22) {
		t.Error("expected port 22 to be denied")
	}
}

func TestIsAllowed_Allow(t *testing.T) {
	tr := portaccess.New()
	_ = tr.Set(80, portaccess.PolicyAllow, "http ok")
	if !tr.IsAllowed(80) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestLen(t *testing.T) {
	tr := portaccess.New()
	if tr.Len() != 0 {
		t.Error("expected empty tracker")
	}
	_ = tr.Set(80, portaccess.PolicyAllow, "")
	_ = tr.Set(443, portaccess.PolicyDeny, "")
	if tr.Len() != 2 {
		t.Errorf("expected 2, got %d", tr.Len())
	}
}

func TestPolicy_String(t *testing.T) {
	if portaccess.PolicyAllow.String() != "allow" {
		t.Error("expected 'allow'")
	}
	if portaccess.PolicyDeny.String() != "deny" {
		t.Error("expected 'deny'")
	}
}
