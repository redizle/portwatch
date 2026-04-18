package portflag_test

import (
	"testing"

	"github.com/user/portwatch/internal/portflag"
)

func TestSet_And_Get(t *testing.T) {
	f := portflag.New()
	if err := f.Set(8080, "suspicious traffic"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fl, ok := f.Get(8080)
	if !ok {
		t.Fatal("expected flag to exist")
	}
	if fl.Reason != "suspicious traffic" {
		t.Errorf("got reason %q, want %q", fl.Reason, "suspicious traffic")
	}
	if fl.FlaggedAt.IsZero() {
		t.Error("FlaggedAt should not be zero")
	}
}

func TestSet_InvalidPort(t *testing.T) {
	f := portflag.New()
	if err := f.Set(0, "test"); err != portflag.ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
	if err := f.Set(70000, "test"); err != portflag.ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestSet_EmptyReason(t *testing.T) {
	f := portflag.New()
	if err := f.Set(443, ""); err != portflag.ErrEmptyReason {
		t.Errorf("expected ErrEmptyReason, got %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	f := portflag.New()
	_, ok := f.Get(9999)
	if ok {
		t.Error("expected no flag for unknown port")
	}
}

func TestIsFlagged(t *testing.T) {
	f := portflag.New()
	_ = f.Set(22, "review")
	if !f.IsFlagged(22) {
		t.Error("expected port 22 to be flagged")
	}
	if f.IsFlagged(80) {
		t.Error("expected port 80 to not be flagged")
	}
}

func TestUnflag(t *testing.T) {
	f := portflag.New()
	_ = f.Set(3306, "db exposed")
	f.Unflag(3306)
	if f.IsFlagged(3306) {
		t.Error("expected port to be unflagged")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	f := portflag.New()
	_ = f.Set(80, "a")
	_ = f.Set(443, "b")
	all := f.All()
	if len(all) != 2 {
		t.Errorf("expected 2 flags, got %d", len(all))
	}
}
