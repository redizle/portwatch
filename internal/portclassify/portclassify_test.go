package portclassify_test

import (
	"testing"

	"portwatch/internal/portclassify"
)

func TestClassify_SystemPort(t *testing.T) {
	c := portclassify.New()
	tier, err := c.Classify(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tier != portclassify.TierSystem {
		t.Errorf("want system, got %s", tier)
	}
}

func TestClassify_RegisteredPort(t *testing.T) {
	c := portclassify.New()
	tier, err := c.Classify(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tier != portclassify.TierRegistered {
		t.Errorf("want registered, got %s", tier)
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	c := portclassify.New()
	tier, err := c.Classify(55000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tier != portclassify.TierDynamic {
		t.Errorf("want dynamic, got %s", tier)
	}
}

func TestClassify_InvalidPort(t *testing.T) {
	c := portclassify.New()
	_, err := c.Classify(0)
	if err == nil {
		t.Error("expected error for port 0")
	}
	_, err = c.Classify(70000)
	if err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestOverride_TakesPrecedence(t *testing.T) {
	c := portclassify.New()
	if err := c.Override(80, portclassify.TierDynamic); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tier, _ := c.Classify(80)
	if tier != portclassify.TierDynamic {
		t.Errorf("want dynamic override, got %s", tier)
	}
}

func TestClearOverride_FallsBackToDefault(t *testing.T) {
	c := portclassify.New()
	_ = c.Override(443, portclassify.TierDynamic)
	c.ClearOverride(443)
	tier, _ := c.Classify(443)
	if tier != portclassify.TierSystem {
		t.Errorf("want system after clear, got %s", tier)
	}
}

func TestOverride_InvalidPort(t *testing.T) {
	c := portclassify.New()
	if err := c.Override(-1, portclassify.TierSystem); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestTier_String(t *testing.T) {
	if portclassify.TierSystem.String() != "system" {
		t.Error("unexpected string for TierSystem")
	}
	if portclassify.TierDynamic.String() != "dynamic" {
		t.Error("unexpected string for TierDynamic")
	}
}
